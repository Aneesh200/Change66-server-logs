package main

import (
	"context"
	"log-ingestion-server/auth"
	"log-ingestion-server/config"
	"log-ingestion-server/database"
	"log-ingestion-server/handlers"
	"log-ingestion-server/middleware"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

const (
	VERSION = "1.0.0"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup logging
	setupLogging(cfg.LogLevel)

	logrus.Infof("Starting Log Ingestion Server v%s", VERSION)
	logrus.Infof("Environment: %s", cfg.Environment)

	// Initialize database
	db, err := database.NewDB(cfg)
	if err != nil {
		logrus.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Run database migrations
	if err := db.RunMigrations(); err != nil {
		logrus.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize authentication service
	authService := auth.NewAuthService(db)
	if err := authService.InitializeAPIKeys(cfg.APIKeys); err != nil {
		logrus.Fatalf("Failed to initialize API keys: %v", err)
	}

	// Initialize handlers
	ingestHandler := handlers.NewIngestHandler(db)
	healthHandler := handlers.NewHealthHandler(db, VERSION)

	// Setup Gin
	gin.SetMode(cfg.GinMode)
	router := gin.New()

	// Apply global middleware
	router.Use(middleware.ErrorHandlingMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.TimeoutMiddleware(cfg.RequestTimeout))
	router.Use(middleware.RequestSizeLimit(cfg.MaxRequestSizeMB))
	
	if cfg.EnableCORS {
		router.Use(middleware.CORSMiddleware(cfg.AllowedOrigins))
	}

	// Rate limiting
	router.Use(middleware.RateLimitMiddleware(cfg.RateLimitRequestsPerMinute, cfg.RateLimitBurst))

	// Public health endpoints (no authentication required)
	router.GET(cfg.HealthCheckPath, healthHandler.HealthCheck)
	router.GET("/readiness", healthHandler.ReadinessCheck)
	router.GET("/liveness", healthHandler.LivenessCheck)

	// Metrics endpoint (if enabled)
	if cfg.EnableMetrics {
		router.GET(cfg.MetricsPath, gin.WrapH(promhttp.Handler()))
	}

	// API v1 routes with authentication
	v1 := router.Group("/api/v1")
	v1.Use(authService.AuthMiddleware())
	{
		// Log ingestion endpoints
		v1.POST("/ingest", ingestHandler.IngestSingle)
		v1.POST("/batch-ingest", ingestHandler.IngestBatch)

		// Admin/monitoring endpoints
		v1.GET("/status", healthHandler.GetStatus)
		v1.GET("/metrics", ingestHandler.GetMetrics)
		v1.GET("/logs/recent", ingestHandler.GetRecentLogs)
		v1.GET("/logs/filter", ingestHandler.GetFilteredLogs)
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:           ":" + cfg.Port,
		Handler:        router,
		ReadTimeout:    cfg.RequestTimeout,
		WriteTimeout:   cfg.RequestTimeout,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// Start server in a goroutine
	go func() {
		logrus.Infof("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Print startup information
	printStartupInfo(cfg)

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Info("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logrus.Errorf("Server forced to shutdown: %v", err)
	}

	logrus.Info("Server exited")
}

// setupLogging configures the logging system
func setupLogging(logLevel string) {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.Warnf("Invalid log level '%s', using 'info'", logLevel)
		level = logrus.InfoLevel
	}

	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	logrus.Info("Logging initialized")
}

// printStartupInfo prints useful startup information
func printStartupInfo(cfg *config.Config) {
	logrus.Info("=== Log Ingestion Server Started ===")
	logrus.Infof("Version: %s", VERSION)
	logrus.Infof("Port: %s", cfg.Port)
	logrus.Infof("Environment: %s", cfg.Environment)
	logrus.Infof("Database: %s:%d/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.Name)
	logrus.Infof("API Keys configured: %d", len(cfg.APIKeys))
	logrus.Infof("Max batch size: %d", cfg.MaxBatchSize)
	logrus.Infof("Worker pool size: %d", cfg.WorkerPoolSize)
	logrus.Infof("Rate limit: %d requests/minute", cfg.RateLimitRequestsPerMinute)
	logrus.Infof("Metrics enabled: %t", cfg.EnableMetrics)
	logrus.Infof("CORS enabled: %t", cfg.EnableCORS)
	
	if cfg.EnableMetrics {
		logrus.Infof("Metrics endpoint: %s", cfg.MetricsPath)
	}
	
	logrus.Info("Available endpoints:")
	logrus.Info("  POST /api/v1/ingest - Single log ingestion")
	logrus.Info("  POST /api/v1/batch-ingest - Batch log ingestion")
	logrus.Info("  GET /health - Health check")
	logrus.Info("  GET /readiness - Readiness check")
	logrus.Info("  GET /liveness - Liveness check")
	logrus.Info("  GET /api/v1/status - Service status")
	logrus.Info("  GET /api/v1/metrics - Analytics metrics")
	logrus.Info("  GET /api/v1/logs/recent - Recent logs")
	logrus.Info("  GET /api/v1/logs/filter - Filtered logs with advanced search")
	
	if cfg.EnableMetrics {
		logrus.Infof("  GET %s - Prometheus metrics", cfg.MetricsPath)
	}
	
	logrus.Info("========================================")
}
