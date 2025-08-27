package auth

import (
	"crypto/sha256"
	"fmt"
	"log-ingestion-server/database"
	"log-ingestion-server/models"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AuthService handles API key authentication
type AuthService struct {
	db       *database.DB
	apiKeys  map[string]bool // In-memory cache for API keys
	lastSync time.Time
}

// NewAuthService creates a new authentication service
func NewAuthService(db *database.DB) *AuthService {
	return &AuthService{
		db:      db,
		apiKeys: make(map[string]bool),
	}
}

// InitializeAPIKeys initializes API keys from configuration
func (as *AuthService) InitializeAPIKeys(keys []string) error {
	for _, key := range keys {
		keyHash := as.hashAPIKey(key)
		
		// Check if API key already exists
		existingKey, err := as.db.GetAPIKey(keyHash)
		if err != nil {
			return fmt.Errorf("failed to check existing API key: %w", err)
		}

		if existingKey == nil {
			// Create new API key
			apiKey := &models.APIKey{
				KeyHash:   keyHash,
				Name:      fmt.Sprintf("Auto-generated key %s", time.Now().Format("2006-01-02")),
				IsActive:  true,
				ExpiresAt: nil, // No expiration for config-based keys
			}

			if err := as.db.InsertAPIKey(apiKey); err != nil {
				return fmt.Errorf("failed to insert API key: %w", err)
			}

			logrus.Infof("Created API key: %s", apiKey.Name)
		}

		// Add to in-memory cache
		as.apiKeys[keyHash] = true
	}

	as.lastSync = time.Now()
	return nil
}

// AuthMiddleware provides API key authentication middleware
func (as *AuthService) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := as.extractAPIKey(c.Request)
		if apiKey == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "API key is required",
			})
			c.Abort()
			return
		}

		keyHash := as.hashAPIKey(apiKey)

		// Check in-memory cache first
		if as.isValidCachedKey(keyHash) {
			// Update usage asynchronously
			go func() {
				if err := as.db.UpdateAPIKeyUsage(keyHash); err != nil {
					logrus.Errorf("Failed to update API key usage: %v", err)
				}
			}()

			c.Set("api_key_hash", keyHash)
			c.Next()
			return
		}

		// Check database if not in cache or cache is stale
		dbKey, err := as.db.GetAPIKey(keyHash)
		if err != nil {
			logrus.Errorf("Failed to validate API key: %v", err)
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "internal_error",
				Message: "Failed to validate API key",
			})
			c.Abort()
			return
		}

		if dbKey == nil || !dbKey.IsActive {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "Invalid or inactive API key",
			})
			c.Abort()
			return
		}

		// Check expiration
		if dbKey.ExpiresAt != nil && dbKey.ExpiresAt.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "API key has expired",
			})
			c.Abort()
			return
		}

		// Update cache and usage
		as.apiKeys[keyHash] = true
		go func() {
			if err := as.db.UpdateAPIKeyUsage(keyHash); err != nil {
				logrus.Errorf("Failed to update API key usage: %v", err)
			}
		}()

		c.Set("api_key_hash", keyHash)
		c.Next()
	}
}

// extractAPIKey extracts API key from request headers
func (as *AuthService) extractAPIKey(r *http.Request) string {
	// Check Authorization header (Bearer token)
	auth := r.Header.Get("Authorization")
	if auth != "" {
		if strings.HasPrefix(auth, "Bearer ") {
			return strings.TrimPrefix(auth, "Bearer ")
		}
		return auth
	}

	// Check X-API-Key header
	return r.Header.Get("X-API-Key")
}

// hashAPIKey creates a SHA-256 hash of the API key
func (as *AuthService) hashAPIKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return fmt.Sprintf("%x", hash)
}

// isValidCachedKey checks if the key is valid in cache
func (as *AuthService) isValidCachedKey(keyHash string) bool {
	// Refresh cache every 5 minutes
	if time.Since(as.lastSync) > 5*time.Minute {
		as.refreshCache()
	}

	return as.apiKeys[keyHash]
}

// refreshCache refreshes the in-memory API key cache
func (as *AuthService) refreshCache() {
	// For now, we keep the existing cache and rely on database validation
	// In a production environment, you might want to periodically refresh from database
	as.lastSync = time.Now()
}

// GenerateAPIKey generates a new API key
func (as *AuthService) GenerateAPIKey(name string, expiresAt *time.Time) (string, error) {
	// Generate a random API key
	key := fmt.Sprintf("hta_%d_%s", time.Now().UnixNano(), generateRandomString(32))
	keyHash := as.hashAPIKey(key)

	apiKey := &models.APIKey{
		KeyHash:   keyHash,
		Name:      name,
		IsActive:  true,
		ExpiresAt: expiresAt,
	}

	if err := as.db.InsertAPIKey(apiKey); err != nil {
		return "", fmt.Errorf("failed to insert API key: %w", err)
	}

	// Add to cache
	as.apiKeys[keyHash] = true

	return key, nil
}

// RevokeAPIKey revokes an API key
func (as *AuthService) RevokeAPIKey(keyHash string) error {
	// Remove from cache
	delete(as.apiKeys, keyHash)

	// Update database (you'd need to implement this method in database package)
	// For now, we'll just remove from cache
	return nil
}

// generateRandomString generates a random string of specified length
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}

// ValidateAPIKeyFormat validates the format of an API key
func ValidateAPIKeyFormat(key string) bool {
	if len(key) < 10 {
		return false
	}
	
	// Basic format validation - you can make this more sophisticated
	return strings.Contains(key, "hta_") || len(key) >= 20
}

// GetAPIKeyInfo returns information about an API key (for admin purposes)
func (as *AuthService) GetAPIKeyInfo(keyHash string) (*models.APIKey, error) {
	return as.db.GetAPIKey(keyHash)
}
