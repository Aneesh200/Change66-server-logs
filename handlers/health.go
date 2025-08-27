package handlers

import (
	"log-ingestion-server/database"
	"log-ingestion-server/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// HealthHandler handles health check and monitoring endpoints
type HealthHandler struct {
	db        *database.DB
	startTime time.Time
	version   string
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *database.DB, version string) *HealthHandler {
	return &HealthHandler{
		db:        db,
		startTime: time.Now(),
		version:   version,
	}
}

// HealthCheck performs a comprehensive health check
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	status := models.HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Version:   h.version,
		Services:  make(map[string]string),
		Uptime:    time.Since(h.startTime).String(),
	}

	// Check database connectivity
	if err := h.db.HealthCheck(); err != nil {
		logrus.Errorf("Database health check failed: %v", err)
		status.Status = "unhealthy"
		status.Services["database"] = "unhealthy"
		
		c.JSON(http.StatusServiceUnavailable, status)
		return
	}
	status.Services["database"] = "healthy"

	// Check database basic operations
	if _, err := h.db.GetLogCount(); err != nil {
		logrus.Errorf("Database query health check failed: %v", err)
		status.Status = "degraded"
		status.Services["database_queries"] = "unhealthy"
	} else {
		status.Services["database_queries"] = "healthy"
	}

	// Determine overall status
	httpStatus := http.StatusOK
	if status.Status == "unhealthy" {
		httpStatus = http.StatusServiceUnavailable
	} else if status.Status == "degraded" {
		httpStatus = http.StatusPartialContent
	}

	c.JSON(httpStatus, status)
}

// ReadinessCheck checks if the service is ready to accept requests
func (h *HealthHandler) ReadinessCheck(c *gin.Context) {
	// Check if database is accessible
	if err := h.db.HealthCheck(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"ready":   false,
			"message": "Database not ready",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ready":   true,
		"message": "Service is ready",
	})
}

// LivenessCheck checks if the service is alive
func (h *HealthHandler) LivenessCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"alive":   true,
		"uptime":  time.Since(h.startTime).String(),
		"version": h.version,
	})
}

// GetStatus returns detailed service status
func (h *HealthHandler) GetStatus(c *gin.Context) {
	logCount, err := h.db.GetLogCount()
	if err != nil {
		logrus.Errorf("Failed to get log count: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve status",
		})
		return
	}

	metrics, err := h.db.GetMetrics()
	if err != nil {
		logrus.Errorf("Failed to get metrics: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "database_error",
			Message: "Failed to retrieve metrics",
		})
		return
	}

	status := gin.H{
		"service": gin.H{
			"name":       "log-ingestion-server",
			"version":    h.version,
			"uptime":     time.Since(h.startTime).String(),
			"started_at": h.startTime.Format(time.RFC3339),
		},
		"database": gin.H{
			"status":     "healthy",
			"total_logs": logCount,
		},
		"metrics": metrics,
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Success: true,
		Message: "Status retrieved successfully",
		Data:    status,
	})
}
