package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/teyk0o/loupi/api/internal/models"
)

// BodyLimit returns a middleware that limits request body size.
// If the Content-Length header exceeds the limit, the request is rejected immediately.
// The body reader is also wrapped to enforce the limit during reading.
func BodyLimit(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxBytes {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, models.ErrorResponse{
				Error:   "payload_too_large",
				Message: "Request body too large",
			})
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}
