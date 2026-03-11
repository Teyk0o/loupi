package models

import (
	"time"

	"github.com/google/uuid"
)

// CustomOption represents a user-configurable option (symptom type, meal category, sport type).
type CustomOption struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"-"`
	Category  string    `json:"category"`
	Value     string    `json:"value"`
	Label     string    `json:"label"`
	Emoji     *string   `json:"emoji,omitempty"`
	SortOrder int       `json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateCustomOptionRequest holds the data needed to create a custom option.
type CreateCustomOptionRequest struct {
	Value string  `json:"value" binding:"required,min=1,max=100"`
	Label string  `json:"label" binding:"required,min=1,max=100"`
	Emoji *string `json:"emoji,omitempty" binding:"omitempty,max=10"`
}

// UpdateCustomOptionRequest holds the data needed to update a custom option.
type UpdateCustomOptionRequest struct {
	Label     *string `json:"label,omitempty" binding:"omitempty,min=1,max=100"`
	Emoji     *string `json:"emoji,omitempty" binding:"omitempty,max=10"`
	SortOrder *int    `json:"sort_order,omitempty" binding:"omitempty,min=0"`
}

// ReorderCustomOptionsRequest holds the ordered list of option IDs for reordering.
type ReorderCustomOptionsRequest struct {
	OrderedIDs []uuid.UUID `json:"ordered_ids" binding:"required,min=1"`
}
