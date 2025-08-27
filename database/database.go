package database

import (
	"context"
	"database/sql"
	"fmt"
	"log-ingestion-server/config"
	"log-ingestion-server/models"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// DB wraps the database connection and provides methods for data operations
type DB struct {
	conn   *sql.DB
	config *config.Config
}

// NewDB creates a new database connection
func NewDB(cfg *config.Config) (*DB, error) {
	conn, err := sql.Open("postgres", cfg.GetDatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	conn.SetMaxOpenConns(cfg.Database.MaxConnections)
	conn.SetMaxIdleConns(cfg.Database.MaxIdleConnections)
	conn.SetConnMaxLifetime(cfg.Database.MaxLifetime)

	// Test connection
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &DB{
		conn:   conn,
		config: cfg,
	}

	logrus.Info("Database connection established successfully")
	return db, nil
}

// RunMigrations runs database migrations
func (db *DB) RunMigrations() error {
	driver, err := postgres.WithInstance(db.conn, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	logrus.Info("Database migrations completed successfully")
	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// InsertLog inserts a single analytics log
func (db *DB) InsertLog(log *models.AnalyticsLog) error {
	query := `
		INSERT INTO analytics_logs (
			event_id, timestamp, event_type, event_name, properties,
			user_id, session_id, app_version, device_info, sequence_number, priority
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at`

	err := db.conn.QueryRow(
		query,
		log.EventID,
		log.Timestamp,
		log.EventType,
		log.EventName,
		log.Properties,
		log.UserID,
		log.SessionID,
		log.AppVersion,
		log.DeviceInfo,
		log.SequenceNumber,
		log.Priority,
	).Scan(&log.ID, &log.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert log: %w", err)
	}

	return nil
}

// InsertLogsBatch inserts multiple analytics logs in a batch
func (db *DB) InsertLogsBatch(logs []models.AnalyticsLog) error {
	if len(logs) == 0 {
		return nil
	}

	tx, err := db.conn.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO analytics_logs (
			event_id, timestamp, event_type, event_name, properties,
			user_id, session_id, app_version, device_info, sequence_number, priority
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, log := range logs {
		_, err = stmt.Exec(
			log.EventID,
			log.Timestamp,
			log.EventType,
			log.EventName,
			log.Properties,
			log.UserID,
			log.SessionID,
			log.AppVersion,
			log.DeviceInfo,
			log.SequenceNumber,
			log.Priority,
		)
		if err != nil {
			return fmt.Errorf("failed to execute batch insert: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetAPIKey retrieves an API key by hash
func (db *DB) GetAPIKey(keyHash string) (*models.APIKey, error) {
	var apiKey models.APIKey
	query := `
		SELECT id, key_hash, name, is_active, created_at, last_used_at, expires_at, usage_count
		FROM api_keys
		WHERE key_hash = $1 AND is_active = true`

	err := db.conn.QueryRow(query, keyHash).Scan(
		&apiKey.ID,
		&apiKey.KeyHash,
		&apiKey.Name,
		&apiKey.IsActive,
		&apiKey.CreatedAt,
		&apiKey.LastUsedAt,
		&apiKey.ExpiresAt,
		&apiKey.UsageCount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // API key not found
		}
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	return &apiKey, nil
}

// UpdateAPIKeyUsage updates the API key usage statistics
func (db *DB) UpdateAPIKeyUsage(keyHash string) error {
	query := `
		UPDATE api_keys
		SET last_used_at = NOW(), usage_count = usage_count + 1
		WHERE key_hash = $1`

	_, err := db.conn.Exec(query, keyHash)
	if err != nil {
		return fmt.Errorf("failed to update API key usage: %w", err)
	}

	return nil
}

// InsertAPIKey inserts a new API key
func (db *DB) InsertAPIKey(apiKey *models.APIKey) error {
	query := `
		INSERT INTO api_keys (key_hash, name, is_active, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	err := db.conn.QueryRow(
		query,
		apiKey.KeyHash,
		apiKey.Name,
		apiKey.IsActive,
		apiKey.ExpiresAt,
	).Scan(&apiKey.ID, &apiKey.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to insert API key: %w", err)
	}

	return nil
}

// GetMetrics retrieves analytics metrics
func (db *DB) GetMetrics() (*models.MetricsResponse, error) {
	metrics := &models.MetricsResponse{}

	// Total logs
	err := db.conn.QueryRow("SELECT COUNT(*) FROM analytics_logs").Scan(&metrics.TotalLogs)
	if err != nil {
		return nil, fmt.Errorf("failed to get total logs: %w", err)
	}

	// Logs in last hour
	err = db.conn.QueryRow(`
		SELECT COUNT(*) FROM analytics_logs 
		WHERE created_at >= NOW() - INTERVAL '1 hour'`).Scan(&metrics.LogsLastHour)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs last hour: %w", err)
	}

	// Logs in last day
	err = db.conn.QueryRow(`
		SELECT COUNT(*) FROM analytics_logs 
		WHERE created_at >= NOW() - INTERVAL '1 day'`).Scan(&metrics.LogsLastDay)
	if err != nil {
		return nil, fmt.Errorf("failed to get logs last day: %w", err)
	}

	// Active sessions (last 30 minutes)
	err = db.conn.QueryRow(`
		SELECT COUNT(DISTINCT session_id) FROM analytics_logs 
		WHERE created_at >= NOW() - INTERVAL '30 minutes' AND session_id IS NOT NULL`).Scan(&metrics.ActiveSessions)
	if err != nil {
		return nil, fmt.Errorf("failed to get active sessions: %w", err)
	}

	// Top event types
	rows, err := db.conn.Query(`
		SELECT event_type, COUNT(*) as count 
		FROM analytics_logs 
		WHERE created_at >= NOW() - INTERVAL '1 day'
		GROUP BY event_type 
		ORDER BY count DESC 
		LIMIT 10`)
	if err != nil {
		return nil, fmt.Errorf("failed to get top event types: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var eventType models.EventTypeCount
		if err := rows.Scan(&eventType.EventType, &eventType.Count); err != nil {
			return nil, fmt.Errorf("failed to scan event type: %w", err)
		}
		metrics.TopEventTypes = append(metrics.TopEventTypes, eventType)
	}

	return metrics, nil
}

// HealthCheck performs a database health check
func (db *DB) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return db.conn.PingContext(ctx)
}

// GetLogCount returns the total number of logs
func (db *DB) GetLogCount() (int64, error) {
	var count int64
	err := db.conn.QueryRow("SELECT COUNT(*) FROM analytics_logs").Scan(&count)
	return count, err
}

// LogFilter represents filtering options for logs
type LogFilter struct {
	EventType     string
	EventName     string
	UserID        string
	SessionID     string
	AppVersion    string
	Priority      string
	ProviderName  string
	StartTime     *time.Time
	EndTime       *time.Time
	Limit         int
	Offset        int
	SortBy        string
	SortOrder     string
}

// GetFilteredLogs returns logs based on filter criteria
func (db *DB) GetFilteredLogs(filter LogFilter) ([]models.AnalyticsLog, int64, error) {
	// Build WHERE clause dynamically
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.EventType != "" {
		conditions = append(conditions, fmt.Sprintf("event_type = $%d", argIndex))
		args = append(args, filter.EventType)
		argIndex++
	}

	if filter.EventName != "" {
		conditions = append(conditions, fmt.Sprintf("event_name = $%d", argIndex))
		args = append(args, filter.EventName)
		argIndex++
	}

	if filter.UserID != "" {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
		args = append(args, filter.UserID)
		argIndex++
	}

	if filter.SessionID != "" {
		conditions = append(conditions, fmt.Sprintf("session_id = $%d", argIndex))
		args = append(args, filter.SessionID)
		argIndex++
	}

	if filter.AppVersion != "" {
		conditions = append(conditions, fmt.Sprintf("app_version = $%d", argIndex))
		args = append(args, filter.AppVersion)
		argIndex++
	}

	if filter.Priority != "" {
		conditions = append(conditions, fmt.Sprintf("priority = $%d", argIndex))
		args = append(args, filter.Priority)
		argIndex++
	}

	if filter.ProviderName != "" {
		conditions = append(conditions, fmt.Sprintf("properties->'tags'->>'provider' = $%d", argIndex))
		args = append(args, filter.ProviderName)
		argIndex++
	}

	if filter.StartTime != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, *filter.StartTime)
		argIndex++
	}

	if filter.EndTime != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, *filter.EndTime)
		argIndex++
	}

	// Build WHERE clause
	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + fmt.Sprintf("%s", conditions[0])
		for i := 1; i < len(conditions); i++ {
			whereClause += " AND " + conditions[i]
		}
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM analytics_logs %s", whereClause)
	var totalCount int64
	err := db.conn.QueryRow(countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get filtered logs count: %w", err)
	}

	// Build ORDER BY clause
	sortBy := "created_at"
	if filter.SortBy != "" {
		// Validate sort field to prevent SQL injection
		validSortFields := map[string]bool{
			"id": true, "event_id": true, "timestamp": true, "event_type": true,
			"event_name": true, "user_id": true, "session_id": true,
			"app_version": true, "priority": true, "created_at": true,
		}
		if validSortFields[filter.SortBy] {
			sortBy = filter.SortBy
		}
	}

	sortOrder := "DESC"
	if filter.SortOrder == "ASC" || filter.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	// Add LIMIT and OFFSET
	limitClause := fmt.Sprintf("ORDER BY %s %s LIMIT $%d", sortBy, sortOrder, argIndex)
	args = append(args, filter.Limit)
	argIndex++

	if filter.Offset > 0 {
		limitClause += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filter.Offset)
	}

	// Build final query
	query := fmt.Sprintf(`
		SELECT id, event_id, timestamp, event_type, event_name, properties,
			   user_id, session_id, app_version, device_info, sequence_number,
			   priority, created_at, processed_at
		FROM analytics_logs 
		%s 
		%s`, whereClause, limitClause)

	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get filtered logs: %w", err)
	}
	defer rows.Close()

	var logs []models.AnalyticsLog
	for rows.Next() {
		var log models.AnalyticsLog
		err := rows.Scan(
			&log.ID,
			&log.EventID,
			&log.Timestamp,
			&log.EventType,
			&log.EventName,
			&log.Properties,
			&log.UserID,
			&log.SessionID,
			&log.AppVersion,
			&log.DeviceInfo,
			&log.SequenceNumber,
			&log.Priority,
			&log.CreatedAt,
			&log.ProcessedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan filtered log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, totalCount, nil
}

// GetRecentLogs returns recent logs for debugging
func (db *DB) GetRecentLogs(limit int) ([]models.AnalyticsLog, error) {
	query := `
		SELECT id, event_id, timestamp, event_type, event_name, properties,
			   user_id, session_id, app_version, device_info, sequence_number,
			   priority, created_at, processed_at
		FROM analytics_logs 
		ORDER BY created_at DESC 
		LIMIT $1`

	rows, err := db.conn.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent logs: %w", err)
	}
	defer rows.Close()

	var logs []models.AnalyticsLog
	for rows.Next() {
		var log models.AnalyticsLog
		err := rows.Scan(
			&log.ID,
			&log.EventID,
			&log.Timestamp,
			&log.EventType,
			&log.EventName,
			&log.Properties,
			&log.UserID,
			&log.SessionID,
			&log.AppVersion,
			&log.DeviceInfo,
			&log.SequenceNumber,
			&log.Priority,
			&log.CreatedAt,
			&log.ProcessedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}
