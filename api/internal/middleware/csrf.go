package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/teyk0o/loupi/api/internal/models"
)

const (
	csrfCookieName = "loupi_csrf"
	csrfHeaderName = "X-CSRF-Token"
	csrfTokenBytes = 32
)

// CSRF implements double-submit cookie CSRF protection.
// On safe methods (GET, HEAD, OPTIONS), a CSRF token cookie is set if not present.
// On state-changing methods (POST, PUT, DELETE), the X-CSRF-Token header must match
// the loupi_csrf cookie value.
func CSRF(cookieDomain string, cookieSecure bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		// Safe methods: ensure CSRF cookie is set
		if method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions {
			if _, err := c.Cookie(csrfCookieName); err != nil {
				token := generateCSRFToken()
				c.SetSameSite(http.SameSiteStrictMode)
				// CSRF cookie must NOT be httpOnly — JavaScript needs to read it
				c.SetCookie(csrfCookieName, token, 86400, "/", cookieDomain, cookieSecure, false)
			}
			c.Next()
			return
		}

		// State-changing methods: validate CSRF token
		cookieToken, err := c.Cookie(csrfCookieName)
		if err != nil || cookieToken == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, models.ErrorResponse{
				Error:   "csrf_error",
				Message: "Missing CSRF token",
			})
			return
		}

		headerToken := c.GetHeader(csrfHeaderName)
		if headerToken == "" || headerToken != cookieToken {
			c.AbortWithStatusJSON(http.StatusForbidden, models.ErrorResponse{
				Error:   "csrf_error",
				Message: "Invalid CSRF token",
			})
			return
		}

		c.Next()
	}
}

func generateCSRFToken() string {
	b := make([]byte, csrfTokenBytes)
	if _, err := rand.Read(b); err != nil {
		// Fallback should never happen, but don't panic
		return hex.EncodeToString(make([]byte, csrfTokenBytes))
	}
	return hex.EncodeToString(b)
}
