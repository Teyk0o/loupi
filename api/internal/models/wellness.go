package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// WellnessEntry represents a daily wellness tracking entry.
type WellnessEntry struct {
	ID           uuid.UUID       `json:"id"`
	UserID       uuid.UUID       `json:"user_id"`
	Date         time.Time       `json:"date"`
	Stress       *int            `json:"stress,omitempty"`
	Mood         *int            `json:"mood,omitempty"`
	Energy       *int            `json:"energy,omitempty"`
	SleepHours   *float32        `json:"sleep_hours,omitempty"`
	SleepQuality *int            `json:"sleep_quality,omitempty"`
	Sport        json.RawMessage `json:"sport,omitempty"`
	Hydration    *int            `json:"hydration,omitempty"`
	Notes        *string         `json:"notes,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// CreateWellnessRequest holds the data needed to create or update a wellness entry.
type CreateWellnessRequest struct {
	Date         string           `json:"date" binding:"required"`
	Stress       *int             `json:"stress,omitempty" binding:"omitempty,min=0,max=5"`
	Mood         *int             `json:"mood,omitempty" binding:"omitempty,min=0,max=5"`
	Energy       *int             `json:"energy,omitempty" binding:"omitempty,min=0,max=5"`
	SleepHours   *float32         `json:"sleep_hours,omitempty" binding:"omitempty,min=0,max=24"`
	SleepQuality *int             `json:"sleep_quality,omitempty" binding:"omitempty,min=0,max=5"`
	Sport        *json.RawMessage `json:"sport,omitempty"`
	Hydration    *int             `json:"hydration,omitempty" binding:"omitempty,min=0,max=50"`
	Notes        *string          `json:"notes,omitempty" binding:"omitempty,max=2000"`
}
