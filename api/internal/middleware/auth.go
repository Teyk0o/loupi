// Package middleware provides HTTP middleware for the Gin router.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/teyk0o/loupi/api/internal/models"
	"github.com/teyk0o/loupi/api/internal/services"
)

// contextKey is the key used to store the user ID in the Gin context.
const contextKeyUserID = "userID"

// Auth returns a middleware that validates JWT access tokens.
// It first checks for the loupi_access cookie, then falls back to the
// Authorization header (Bearer scheme) for API clients.
// Tokens are also checked against the Redis blacklist.
func Auth(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// Try cookie first (browser clients)
		if cookie, err := c.Cookie("loupi_access"); err == nil && cookie != "" {
			tokenString = cookie
		}

		// Fall back to Authorization header (API clients)
		if tokenString == "" {
			header := c.GetHeader("Authorization")
			if header != "" {
				parts := strings.SplitN(header, " ", 2)
				if len(parts) == 2 && parts[0] == "Bearer" {
					tokenString = parts[1]
				}
			}
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "Authentication required",
			})
			return
		}

		userID, _, err := authService.ValidateAccessToken(c.Request.Context(), tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "unauthorized",
				Message: "Invalid or expired token",
			})
			return
		}

		c.Set(contextKeyUserID, userID)
		c.Next()
	}
}

// GetUserID extracts the authenticated user ID from the Gin context.
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get(contextKeyUserID)
	if !exists {
		return uuid.Nil, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}
