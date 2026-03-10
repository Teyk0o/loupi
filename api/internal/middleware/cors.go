package middleware

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS returns a middleware that handles Cross-Origin Resource Sharing.
// Wildcard origins are rejected — the server will refuse to start if "*" is configured.
func CORS(allowedOrigins string) gin.HandlerFunc {
	origins := make([]string, 0)
	for _, o := range strings.Split(allowedOrigins, ",") {
		trimmed := strings.TrimSpace(o)
		if trimmed == "" {
			continue
		}
		if trimmed == "*" {
			log.Fatal("CORS: wildcard '*' origin is not allowed with credentials. Set explicit origins in LOUPI_ALLOWED_ORIGINS.")
		}
		origins = append(origins, trimmed)
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		c.Header("Vary", "Origin")

		matched := false
		if origin != "" {
			for _, allowed := range origins {
				if allowed == origin {
					matched = true
					break
				}
			}
		}

		if matched {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type, X-CSRF-Token")
			c.Header("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
