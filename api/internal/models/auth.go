package models

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

// RegisterPasswordComplexity is a custom validator for password complexity.
// Requires at least 1 uppercase, 1 lowercase, 1 digit, and 1 special character.
func RegisterPasswordComplexity(v *validator.Validate) error {
	return v.RegisterValidation("password_complexity", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		var hasUpper, hasLower, hasDigit, hasSpecial bool
		for _, ch := range password {
			switch {
			case unicode.IsUpper(ch):
				hasUpper = true
			case unicode.IsLower(ch):
				hasLower = true
			case unicode.IsDigit(ch):
				hasDigit = true
			case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
				hasSpecial = true
			}
		}
		return hasUpper && hasLower && hasDigit && hasSpecial
	})
}

// RegisterRequest holds the data needed to register a new user.
type RegisterRequest struct {
	Email     string  `json:"email" binding:"required,email,max=255"`
	Password  string  `json:"password" binding:"required,min=8,max=128,password_complexity"`
	FirstName *string `json:"first_name,omitempty" binding:"omitempty,max=100"`
}

// LoginRequest holds the data needed to authenticate a user.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse is the internal token response used for cookie setting.
type AuthResponse struct {
	AccessToken  string       `json:"-"`
	RefreshToken string       `json:"-"`
	ExpiresIn    int64        `json:"expires_in"`
	User         UserResponse `json:"user"`
}

// CookieAuthResponse is the public response after login/register (no tokens in body).
type CookieAuthResponse struct {
	ExpiresIn int64        `json:"expires_in"`
	User      UserResponse `json:"user"`
}

// RefreshRequest is kept for backward compatibility.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// ForgotPasswordRequest holds the data needed to request a password reset.
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest holds the data needed to reset a password.
type ResetPasswordRequest struct {
	Token    string `json:"token" binding:"required"`
	Password string `json:"password" binding:"required,min=8,max=128,password_complexity"`
}
