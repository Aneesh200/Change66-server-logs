package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// JSONB represents a PostgreSQL JSONB field
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface for JSONB
func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for JSONB
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into JSONB", value)
	}

	return json.Unmarshal(bytes, j)
}

// AnalyticsLog represents a single analytics log entry
type AnalyticsLog struct {
	ID             int64     `json:"id" db:"id"`
	EventID        string    `json:"event_id" db:"event_id" validate:"required"`
	Timestamp      time.Time `json:"timestamp" db:"timestamp" validate:"required"`
	EventType      string    `json:"event_type" db:"event_type" validate:"required,oneof=behavioral telemetry observability error performance"`
	EventName      string    `json:"event_name" db:"event_name" validate:"required"`
	Properties     JSONB     `json:"properties" db:"properties"`
	UserID         *string   `json:"user_id" db:"user_id"`
	SessionID      *string   `json:"session_id" db:"session_id"`
	AppVersion     *string   `json:"app_version" db:"app_version"`
	DeviceInfo     JSONB     `json:"device_info" db:"device_info"`
	SequenceNumber *int      `json:"sequence_number" db:"sequence_number"`
	Priority       string    `json:"priority" db:"priority" validate:"oneof=normal high"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	ProcessedAt    *time.Time `json:"processed_at" db:"processed_at"`
}

// BatchRequest represents a batch of analytics logs
type BatchRequest struct {
	Logs []AnalyticsLog `json:"logs" validate:"required,dive"`
}

// APIKey represents an API key in the database
type APIKey struct {
	ID          int64      `json:"id" db:"id"`
	KeyHash     string     `json:"-" db:"key_hash"`
	Name        string     `json:"name" db:"name"`
	IsActive    bool       `json:"is_active" db:"is_active"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	LastUsedAt  *time.Time `json:"last_used_at" db:"last_used_at"`
	ExpiresAt   *time.Time `json:"expires_at" db:"expires_at"`
	UsageCount  int64      `json:"usage_count" db:"usage_count"`
}

// ServerMetric represents a server metric entry
type ServerMetric struct {
	ID          int64     `json:"id" db:"id"`
	MetricName  string    `json:"metric_name" db:"metric_name"`
	MetricValue float64   `json:"metric_value" db:"metric_value"`
	Labels      JSONB     `json:"labels" db:"labels"`
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
}

// ValidationError represents a validation error response
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Details []ValidationError `json:"details,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// HealthStatus represents the health check response
type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	Services  map[string]string `json:"services"`
	Uptime    string            `json:"uptime"`
}

// MetricsResponse represents the metrics response
type MetricsResponse struct {
	TotalLogs       int64   `json:"total_logs"`
	LogsLastHour    int64   `json:"logs_last_hour"`
	LogsLastDay     int64   `json:"logs_last_day"`
	AverageLatency  float64 `json:"average_latency_ms"`
	ErrorRate       float64 `json:"error_rate_percent"`
	ActiveSessions  int64   `json:"active_sessions"`
	TopEventTypes   []EventTypeCount `json:"top_event_types"`
}

// EventTypeCount represents event type counts
type EventTypeCount struct {
	EventType string `json:"event_type"`
	Count     int64  `json:"count"`
}

// FilteredLogsResponse represents the response for filtered logs
type FilteredLogsResponse struct {
	Logs       []AnalyticsLog `json:"logs"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalPages int            `json:"total_pages"`
}
