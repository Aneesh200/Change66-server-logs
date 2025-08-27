package handlers

import (
	"fmt"
	"log-ingestion-server/database"
	"log-ingestion-server/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

// IngestHandler handles log ingestion requests
type IngestHandler struct {
	db        *database.DB
	validator *validator.Validate
	metrics   *Metrics
}

// Metrics holds Prometheus metrics
type Metrics struct {
	RequestsTotal     *prometheus.CounterVec
	RequestDuration   *prometheus.HistogramVec
	LogsIngested      *prometheus.CounterVec
	BatchSize         *prometheus.HistogramVec
	ValidationErrors  *prometheus.CounterVec
	DatabaseErrors    *prometheus.CounterVec
}

// NewIngestHandler creates a new ingest handler
func NewIngestHandler(db *database.DB) *IngestHandler {
	validator := validator.New()
	
	// Register custom validation for event types
	validator.RegisterValidation("event_type", validateEventType)

	metrics := &Metrics{
		RequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status"},
		),
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "http_request_duration_seconds",
				Help: "HTTP request duration in seconds",
			},
			[]string{"method", "endpoint"},
		),
		LogsIngested: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "logs_ingested_total",
				Help: "Total number of logs ingested",
			},
			[]string{"event_type", "priority"},
		),
		BatchSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "batch_size",
				Help: "Size of log batches",
				Buckets: []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000},
			},
			[]string{"endpoint"},
		),
		ValidationErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "validation_errors_total",
				Help: "Total number of validation errors",
			},
			[]string{"field", "error_type"},
		),
		DatabaseErrors: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "database_errors_total",
				Help: "Total number of database errors",
			},
			[]string{"operation"},
		),
	}

	// Register metrics
	prometheus.MustRegister(
		metrics.RequestsTotal,
		metrics.RequestDuration,
		metrics.LogsIngested,
		metrics.BatchSize,
		metrics.ValidationErrors,
		metrics.DatabaseErrors,
	)

	return &IngestHandler{
		db:        db,
		validator: validator,
		metrics:   metrics,
	}
}

// IngestSingle handles single log ingestion
func (h *IngestHandler) IngestSingle(c *gin.Context) {
	start := time.Now()
	
	defer func() {
		duration := time.Since(start).Seconds()
		h.metrics.RequestDuration.WithLabelValues("POST", "/ingest").Observe(duration)
		h.metrics.RequestsTotal.WithLabelValues("POST", "/ingest", fmt.Sprintf("%d", c.Writer.Status())).Inc()
	}()

	var log models.AnalyticsLog
	if err := c.ShouldBindJSON(&log); err != nil {
		logrus.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_json",
			Message: "Invalid JSON format",
		})
		return
	}

	// Set default values
	h.setDefaultValues(&log)

	// Validate the log
	if err := h.validator.Struct(&log); err != nil {
		validationErrors := h.formatValidationErrors(err)
		for _, ve := range validationErrors {
			h.metrics.ValidationErrors.WithLabelValues(ve.Field, "validation").Inc()
		}
		
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid log data",
			Details: validationErrors,
		})
		return
	}

	// Insert into database
	if err := h.db.InsertLog(&log); err != nil {
		logrus.Errorf("Failed to insert log: %v", err)
		h.metrics.DatabaseErrors.WithLabelValues("insert_single").Inc()
		
		if strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Error:   "duplicate_event",
				Message: "Event with this ID already exists",
			})
			return
		}
		
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to store log",
		})
		return
	}

	// Update metrics
	h.metrics.LogsIngested.WithLabelValues(log.EventType, log.Priority).Inc()
	h.metrics.BatchSize.WithLabelValues("single").Observe(1)

	logrus.Debugf("Ingested single log: %s", log.EventID)

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Success: true,
		Message: "Log ingested successfully",
		Data: map[string]interface{}{
			"event_id": log.EventID,
			"id":       log.ID,
		},
	})
}

// IngestBatch handles batch log ingestion
func (h *IngestHandler) IngestBatch(c *gin.Context) {
	start := time.Now()
	
	defer func() {
		duration := time.Since(start).Seconds()
		h.metrics.RequestDuration.WithLabelValues("POST", "/batch-ingest").Observe(duration)
		h.metrics.RequestsTotal.WithLabelValues("POST", "/batch-ingest", fmt.Sprintf("%d", c.Writer.Status())).Inc()
	}()

	var batchRequest models.BatchRequest
	if err := c.ShouldBindJSON(&batchRequest); err != nil {
		logrus.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_json",
			Message: "Invalid JSON format",
		})
		return
	}

	if len(batchRequest.Logs) == 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "empty_batch",
			Message: "Batch cannot be empty",
		})
		return
	}

	// Validate batch size
	maxBatchSize := 1000 // This should come from config
	if len(batchRequest.Logs) > maxBatchSize {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "batch_too_large",
			Message: fmt.Sprintf("Batch size cannot exceed %d logs", maxBatchSize),
		})
		return
	}

	// Set default values and validate each log
	var validationErrors []models.ValidationError
	validLogs := make([]models.AnalyticsLog, 0, len(batchRequest.Logs))

	for i, log := range batchRequest.Logs {
		h.setDefaultValues(&log)
		
		if err := h.validator.Struct(&log); err != nil {
			for _, ve := range h.formatValidationErrors(err) {
				ve.Field = fmt.Sprintf("logs[%d].%s", i, ve.Field)
				validationErrors = append(validationErrors, ve)
				h.metrics.ValidationErrors.WithLabelValues(ve.Field, "validation").Inc()
			}
		} else {
			validLogs = append(validLogs, log)
		}
	}

	// If there are validation errors, return them
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: fmt.Sprintf("Validation failed for %d logs", len(validationErrors)),
			Details: validationErrors,
		})
		return
	}

	// Insert batch into database
	if err := h.db.InsertLogsBatch(validLogs); err != nil {
		logrus.Errorf("Failed to insert batch: %v", err)
		h.metrics.DatabaseErrors.WithLabelValues("insert_batch").Inc()
		
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to store logs",
		})
		return
	}

	// Update metrics
	for _, log := range validLogs {
		h.metrics.LogsIngested.WithLabelValues(log.EventType, log.Priority).Inc()
	}
	h.metrics.BatchSize.WithLabelValues("batch").Observe(float64(len(validLogs)))

	logrus.Infof("Ingested batch of %d logs", len(validLogs))

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Success: true,
		Message: fmt.Sprintf("Batch of %d logs ingested successfully", len(validLogs)),
		Data: map[string]interface{}{
			"logs_processed": len(validLogs),
			"total_received": len(batchRequest.Logs),
		},
	})
}

// setDefaultValues sets default values for log fields
func (h *IngestHandler) setDefaultValues(log *models.AnalyticsLog) {
	if log.Timestamp.IsZero() {
		log.Timestamp = time.Now()
	}
	
	if log.Priority == "" {
		log.Priority = "normal"
	}
	
	if log.Properties == nil {
		log.Properties = make(models.JSONB)
	}
	
	if log.DeviceInfo == nil {
		log.DeviceInfo = make(models.JSONB)
	}
}

// formatValidationErrors formats validation errors into a readable format
func (h *IngestHandler) formatValidationErrors(err error) []models.ValidationError {
	var validationErrors []models.ValidationError
	
	if validatorErr, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validatorErr {
			validationErrors = append(validationErrors, models.ValidationError{
				Field:   fieldError.Field(),
				Message: getValidationMessage(fieldError),
			})
		}
	}
	
	return validationErrors
}

// getValidationMessage returns a human-readable validation message
func getValidationMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "This field is required"
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", fieldError.Param())
	case "min":
		return fmt.Sprintf("Must be at least %s characters", fieldError.Param())
	case "max":
		return fmt.Sprintf("Must be at most %s characters", fieldError.Param())
	case "email":
		return "Must be a valid email address"
	case "event_type":
		return "Must be a valid event type (behavioral, telemetry, observability, error, performance)"
	default:
		return fmt.Sprintf("Validation failed for tag '%s'", fieldError.Tag())
	}
}

// validateEventType is a custom validation function for event types
func validateEventType(fl validator.FieldLevel) bool {
	eventType := fl.Field().String()
	validTypes := []string{"behavioral", "telemetry", "observability", "error", "performance"}
	
	for _, validType := range validTypes {
		if eventType == validType {
			return true
		}
	}
	
	return false
}

// GetMetrics returns current metrics (for admin/monitoring purposes)
func (h *IngestHandler) GetMetrics(c *gin.Context) {
	metrics, err := h.db.GetMetrics()
	if err != nil {
		logrus.Errorf("Failed to get metrics: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve metrics",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Metrics retrieved successfully",
		Data:    metrics,
	})
}

// GetRecentLogs returns recent logs for debugging
func (h *IngestHandler) GetRecentLogs(c *gin.Context) {
	limit := 50 // Default limit
	if limitParam := c.Query("limit"); limitParam != "" {
		if parsedLimit, err := parseInt(limitParam); err == nil && parsedLimit > 0 && parsedLimit <= 1000 {
			limit = parsedLimit
		}
	}

	logs, err := h.db.GetRecentLogs(limit)
	if err != nil {
		logrus.Errorf("Failed to get recent logs: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve logs",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: fmt.Sprintf("Retrieved %d recent logs", len(logs)),
		Data:    logs,
	})
}

// GetFilteredLogs returns filtered logs based on query parameters
func (h *IngestHandler) GetFilteredLogs(c *gin.Context) {
	start := time.Now()
	
	defer func() {
		duration := time.Since(start).Seconds()
		h.metrics.RequestDuration.WithLabelValues("GET", "/logs/filter").Observe(duration)
		h.metrics.RequestsTotal.WithLabelValues("GET", "/logs/filter", fmt.Sprintf("%d", c.Writer.Status())).Inc()
	}()

	// Parse query parameters
	filter := database.LogFilter{
		EventType:    c.Query("event_type"),
		EventName:    c.Query("event_name"),
		UserID:       c.Query("user_id"),
		SessionID:    c.Query("session_id"),
		AppVersion:   c.Query("app_version"),
		Priority:     c.Query("priority"),
		ProviderName: c.Query("provider_name"),
		SortBy:       c.Query("sort_by"),
		SortOrder:    c.Query("sort_order"),
	}

	// Parse pagination parameters
	page := 1
	if pageParam := c.Query("page"); pageParam != "" {
		if parsedPage, err := strconv.Atoi(pageParam); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	pageSize := 50 // Default page size
	if pageSizeParam := c.Query("page_size"); pageSizeParam != "" {
		if parsedPageSize, err := strconv.Atoi(pageSizeParam); err == nil && parsedPageSize > 0 && parsedPageSize <= 1000 {
			pageSize = parsedPageSize
		}
	}

	filter.Limit = pageSize
	filter.Offset = (page - 1) * pageSize

	// Parse time filters
	if startTimeParam := c.Query("start_time"); startTimeParam != "" {
		if startTime, err := time.Parse(time.RFC3339, startTimeParam); err == nil {
			filter.StartTime = &startTime
		} else {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "invalid_start_time",
				Message: "start_time must be in RFC3339 format (e.g., 2023-01-01T00:00:00Z)",
			})
			return
		}
	}

	if endTimeParam := c.Query("end_time"); endTimeParam != "" {
		if endTime, err := time.Parse(time.RFC3339, endTimeParam); err == nil {
			filter.EndTime = &endTime
		} else {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "invalid_end_time",
				Message: "end_time must be in RFC3339 format (e.g., 2023-01-01T00:00:00Z)",
			})
			return
		}
	}

	// Validate event type if provided
	if filter.EventType != "" {
		validTypes := []string{"behavioral", "telemetry", "observability", "error", "performance"}
		valid := false
		for _, validType := range validTypes {
			if filter.EventType == validType {
				valid = true
				break
			}
		}
		if !valid {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error:   "invalid_event_type",
				Message: "event_type must be one of: behavioral, telemetry, observability, error, performance",
			})
			return
		}
	}

	// Validate priority if provided
	if filter.Priority != "" && filter.Priority != "normal" && filter.Priority != "high" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid_priority",
			Message: "priority must be either 'normal' or 'high'",
		})
		return
	}

	// Get filtered logs
	logs, totalCount, err := h.db.GetFilteredLogs(filter)
	if err != nil {
		logrus.Errorf("Failed to get filtered logs: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve filtered logs",
		})
		return
	}

	// Calculate total pages
	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize != 0 {
		totalPages++
	}

	response := models.FilteredLogsResponse{
		Logs:       logs,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: fmt.Sprintf("Retrieved %d filtered logs (page %d of %d)", len(logs), page, totalPages),
		Data:    response,
	})
}

// parseInt safely parses an integer from string
func parseInt(s string) (int, error) {
	var result int
	if _, err := fmt.Sscanf(s, "%d", &result); err != nil {
		return 0, err
	}
	return result, nil
}
