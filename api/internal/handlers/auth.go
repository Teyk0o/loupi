// Package handlers contains the HTTP request handlers for the API.
package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/teyk0o/loupi/api/internal/config"
	"github.com/teyk0o/loupi/api/internal/middleware"
	"github.com/teyk0o/loupi/api/internal/models"
	"github.com/teyk0o/loupi/api/internal/services"
	"github.com/teyk0o/loupi/api/internal/utils"
)

const (
	accessCookieName  = "loupi_access"
	refreshCookieName = "loupi_refresh"
	accessMaxAge      = 900    // 15 minutes
	refreshMaxAge     = 604800 // 7 days
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authService  *services.AuthService
	audit        *utils.AuditLogger
	cfg          *config.Config
	loginLimiter *middleware.LoginRateLimiter
}

// NewAuthHandler creates a new AuthHandler instance.
func NewAuthHandler(authService *services.AuthService, audit *utils.AuditLogger, cfg *config.Config, loginLimiter *middleware.LoginRateLimiter) *AuthHandler {
	return &AuthHandler{authService: authService, audit: audit, cfg: cfg, loginLimiter: loginLimiter}
}

// setAuthCookies sets httpOnly secure cookies for access and refresh tokens.
func (h *AuthHandler) setAuthCookies(c *gin.Context, tokens *models.AuthResponse) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(accessCookieName, tokens.AccessToken, accessMaxAge, "/", h.cfg.CookieDomain, h.cfg.CookieSecure, true)
	c.SetCookie(refreshCookieName, tokens.RefreshToken, refreshMaxAge, "/v1/auth", h.cfg.CookieDomain, h.cfg.CookieSecure, true)
}

// clearAuthCookies removes auth cookies.
func (h *AuthHandler) clearAuthCookies(c *gin.Context) {
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(accessCookieName, "", -1, "/", h.cfg.CookieDomain, h.cfg.CookieSecure, true)
	c.SetCookie(refreshCookieName, "", -1, "/v1/auth", h.cfg.CookieDomain, h.cfg.CookieSecure, true)
}

// Register handles user registration (POST /v1/auth/register).
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid input data",
		})
		return
	}

	tokens, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, models.ErrorResponse{
				Error:   "conflict",
				Message: "Unable to create account with this email",
			})
			return
		}
		log.Printf("register error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to create account",
		})
		return
	}

	h.setAuthCookies(c, tokens)
	_ = h.audit.Log(c.Request.Context(), tokens.User.ID, "register", "user", tokens.User.ID, c.ClientIP())

	c.JSON(http.StatusCreated, models.CookieAuthResponse{
		ExpiresIn: tokens.ExpiresIn,
		User:      tokens.User,
	})
}

// Login handles user authentication (POST /v1/auth/login).
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "Invalid input data",
		})
		return
	}

	// Check account-level lockout
	if h.loginLimiter.IsLocked(req.Email) {
		c.JSON(http.StatusTooManyRequests, models.ErrorResponse{
			Error:   "account_locked",
			Message: "Too many failed attempts, please try again later",
		})
		return
	}

	tokens, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			h.loginLimiter.RecordFailure(req.Email)
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error:   "invalid_credentials",
				Message: "Invalid email or password",
			})
			return
		}
		log.Printf("login error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to authenticate",
		})
		return
	}

	h.loginLimiter.RecordSuccess(req.Email)
	h.setAuthCookies(c, tokens)
	_ = h.audit.Log(c.Request.Context(), tokens.User.ID, "login", "user", tokens.User.ID, c.ClientIP())

	c.JSON(http.StatusOK, models.CookieAuthResponse{
		ExpiresIn: tokens.ExpiresIn,
		User:      tokens.User,
	})
}

// Refresh handles token refresh (POST /v1/auth/refresh).
// Reads the refresh token from the httpOnly cookie.
func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(refreshCookieName)
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "invalid_token",
			Message: "Missing refresh token",
		})
		return
	}

	tokens, err := h.authService.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		h.clearAuthCookies(c)
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "invalid_token",
			Message: "Invalid or expired refresh token",
		})
		return
	}

	h.setAuthCookies(c, tokens)

	c.JSON(http.StatusOK, models.CookieAuthResponse{
		ExpiresIn: tokens.ExpiresIn,
		User:      tokens.User,
	})
}

// Logout handles user logout (POST /v1/auth/logout).
// Revokes the refresh token and blacklists the access token.
func (h *AuthHandler) Logout(c *gin.Context) {
	accessToken, _ := c.Cookie(accessCookieName)
	refreshToken, _ := c.Cookie(refreshCookieName)

	_ = h.authService.Logout(c.Request.Context(), refreshToken, accessToken)

	userID, ok := middleware.GetUserID(c)
	if ok {
		_ = h.audit.Log(c.Request.Context(), userID, "logout", "user", userID, c.ClientIP())
	}

	h.clearAuthCookies(c)
	c.JSON(http.StatusOK, models.MessageResponse{
		Message: "Logged out successfully",
	})
}

// Me returns the authenticated user's profile (GET /v1/auth/me).
func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "Not authenticated",
		})
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error:   "not_found",
			Message: "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// DeleteAccount handles account deletion (DELETE /v1/auth/account).
func (h *AuthHandler) DeleteAccount(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Error:   "unauthorized",
			Message: "Not authenticated",
		})
		return
	}

	_ = h.audit.Log(c.Request.Context(), userID, "delete_account", "user", userID, c.ClientIP())

	if err := h.authService.DeleteAccount(c.Request.Context(), userID); err != nil {
		log.Printf("delete account error: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to delete account",
		})
		return
	}

	h.clearAuthCookies(c)
	c.JSON(http.StatusOK, models.MessageResponse{
		Message: "Account deleted successfully",
	})
}
