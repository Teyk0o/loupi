// Package models defines the data structures used across the application.
package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a registered user account.
type User struct {
	ID            uuid.UUID  `json:"id"`
	Email         string     `json:"email"`
	PasswordHash  *string    `json:"-"`
	OAuthProvider *string    `json:"oauth_provider,omitempty"`
	OAuthID       *string    `json:"-"`
	EmailVerified bool       `json:"email_verified"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// UserResponse is the public representation of a user (no sensitive fields).
type UserResponse struct {
	ID            uuid.UUID `json:"id"`
	Email         string    `json:"email"`
	EmailVerified bool      `json:"email_verified"`
	CreatedAt     time.Time `json:"created_at"`
}

// ToResponse converts a User to its public representation.
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:            u.ID,
		Email:         u.Email,
		EmailVerified: u.EmailVerified,
		CreatedAt:     u.CreatedAt,
	}
}
