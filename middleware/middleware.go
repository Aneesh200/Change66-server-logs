package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log-ingestion-server/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

// RequestSizeLimit creates middleware to limit request body size
func RequestSizeLimit(maxSizeMB int) gin.HandlerFunc {
	maxSize := int64(maxSizeMB * 1024 * 1024) // Convert MB to bytes
	
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			c.JSON(http.StatusRequestEntityTooLarge, models.ErrorResponse{
				Error:   "request_too_large",
				Message: fmt.Sprintf("Request body too large. Maximum size is %d MB", maxSizeMB),
			})
			c.Abort()
			return
		}
		
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	}
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(requestsPerMinute int, burst int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(requestsPerMinute)/60, burst) // Convert per minute to per second
	
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, models.ErrorResponse{
				Error:   "rate_limit_exceeded",
				Message: "Too many requests. Please slow down.",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// LoggingMiddleware creates a structured logging middleware
func LoggingMiddleware() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logEntry := logrus.WithFields(logrus.Fields{
			"method":     param.Method,
			"path":       param.Path,
			"status":     param.StatusCode,
			"latency":    param.Latency,
			"client_ip":  param.ClientIP,
			"user_agent": param.Request.UserAgent(),
			"error":      param.ErrorMessage,
		})

		if param.StatusCode >= 400 {
			logEntry.Error("HTTP request completed with error")
		} else {
			logEntry.Info("HTTP request completed")
		}

		return ""
	})
}

// ErrorHandlingMiddleware provides centralized error handling
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logrus.Errorf("Panic recovered: %v", recovered)
		
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_server_error",
			Message: "An unexpected error occurred",
		})
	})
}

// SecurityHeadersMiddleware adds security headers
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}

// CORSMiddleware creates CORS middleware with custom configuration
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	config := cors.DefaultConfig()
	
	if len(allowedOrigins) == 1 && allowedOrigins[0] == "*" {
		config.AllowAllOrigins = true
	} else {
		config.AllowOrigins = allowedOrigins
	}
	
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-API-Key"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	
	return cors.New(config)
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		c.Next()
	}
}

// TimeoutMiddleware adds request timeout
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Create a context with timeout
		ctx := c.Request.Context()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// RequestValidationMiddleware validates JSON requests
func RequestValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				c.JSON(http.StatusUnsupportedMediaType, models.ErrorResponse{
					Error:   "unsupported_media_type",
					Message: "Content-Type must be application/json",
				})
				c.Abort()
				return
			}

			// Read and validate JSON
			if c.Request.Body != nil {
				bodyBytes, err := io.ReadAll(c.Request.Body)
				if err != nil {
					c.JSON(http.StatusBadRequest, models.ErrorResponse{
						Error:   "invalid_request_body",
						Message: "Could not read request body",
					})
					c.Abort()
					return
				}

				// Restore the body for further processing
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

				// Basic JSON validation
				if len(bodyBytes) == 0 {
					c.JSON(http.StatusBadRequest, models.ErrorResponse{
						Error:   "empty_request_body",
						Message: "Request body cannot be empty",
					})
					c.Abort()
					return
				}
			}
		}
		
		c.Next()
	}
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}

// CompressionMiddleware adds response compression
func CompressionMiddleware() gin.HandlerFunc {
	// This would typically use a compression library like gzip
	// For now, we'll just return a placeholder
	return func(c *gin.Context) {
		c.Next()
	}
}
