package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// Config holds all configuration for the application
type Config struct {
	// Server Configuration
	Port        string
	GinMode     string
	LogLevel    string
	Environment string

	// Database Configuration
	Database DatabaseConfig

	// API Security
	APIKeys                    []string
	EnableAPIKeyRotation       bool
	APIKeyRotationInterval     time.Duration

	// Rate Limiting
	RateLimitRequestsPerMinute int
	RateLimitBurst             int

	// Batch Processing
	MaxBatchSize     int
	BatchTimeout     time.Duration
	WorkerPoolSize   int

	// Monitoring
	EnableMetrics     bool
	MetricsPath       string
	HealthCheckPath   string

	// Security
	EnableCORS           bool
	AllowedOrigins       []string
	RequestTimeout       time.Duration
	MaxRequestSizeMB     int
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host               string
	Port               int
	Name               string
	User               string
	Password           string
	SSLMode            string
	MaxConnections     int
	MaxIdleConnections int
	MaxLifetime        time.Duration
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		logrus.Info("No .env file found, using environment variables")
	}

	config := &Config{
		// Default values
		Port:        getEnv("PORT", "8080"),
		GinMode:     getEnv("GIN_MODE", "release"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Environment: getEnv("ENVIRONMENT", "production"),

		Database: DatabaseConfig{
			Host:               getEnv("DB_HOST", "localhost"),
			Port:               getEnvAsInt("DB_PORT", 5432),
			Name:               getEnv("DB_NAME", "analytics_logs"),
			User:               getEnv("DB_USER", "postgres"),
			Password:           getEnv("DB_PASSWORD", ""),
			SSLMode:            getEnv("DB_SSL_MODE", "disable"),
			MaxConnections:     getEnvAsInt("DB_MAX_CONNECTIONS", 100),
			MaxIdleConnections: getEnvAsInt("DB_MAX_IDLE_CONNECTIONS", 10),
			MaxLifetime:        time.Duration(getEnvAsInt("DB_MAX_LIFETIME_MINUTES", 60)) * time.Minute,
		},

		APIKeys:                    getEnvAsSlice("API_KEYS", ","),
		EnableAPIKeyRotation:       getEnvAsBool("ENABLE_API_KEY_ROTATION", false),
		APIKeyRotationInterval:     time.Duration(getEnvAsInt("API_KEY_ROTATION_INTERVAL_HOURS", 24)) * time.Hour,

		RateLimitRequestsPerMinute: getEnvAsInt("RATE_LIMIT_REQUESTS_PER_MINUTE", 1000),
		RateLimitBurst:             getEnvAsInt("RATE_LIMIT_BURST", 100),

		MaxBatchSize:     getEnvAsInt("MAX_BATCH_SIZE", 1000),
		BatchTimeout:     time.Duration(getEnvAsInt("BATCH_TIMEOUT_SECONDS", 30)) * time.Second,
		WorkerPoolSize:   getEnvAsInt("WORKER_POOL_SIZE", 10),

		EnableMetrics:     getEnvAsBool("ENABLE_METRICS", true),
		MetricsPath:       getEnv("METRICS_PATH", "/metrics"),
		HealthCheckPath:   getEnv("HEALTH_CHECK_PATH", "/health"),

		EnableCORS:           getEnvAsBool("ENABLE_CORS", true),
		AllowedOrigins:       getEnvAsSlice("ALLOWED_ORIGINS", ","),
		RequestTimeout:       time.Duration(getEnvAsInt("REQUEST_TIMEOUT_SECONDS", 30)) * time.Second,
		MaxRequestSizeMB:     getEnvAsInt("MAX_REQUEST_SIZE_MB", 10),
	}

	// Validate required configuration
	if len(config.APIKeys) == 0 {
		return nil, fmt.Errorf("API_KEYS must be provided")
	}

	if config.Database.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD must be provided")
	}

	return config, nil
}

// GetDatabaseURL returns the database connection string
func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

// Helper functions for environment variables
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getEnvAsSlice(key, delimiter string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, delimiter)
	}
	return []string{}
}
