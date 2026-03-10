package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// SymptomDetail represents a single symptom with its severity.
type SymptomDetail struct {
	Type     string `json:"type" binding:"required,max=100"`
	Severity int    `json:"severity" binding:"required,min=1,max=5"`
}

// SymptomCheckin represents a symptom check-in linked to a meal.
type SymptomCheckin struct {
	ID         uuid.UUID       `json:"id"`
	MealID     uuid.UUID       `json:"meal_id"`
	DelayHours int             `json:"delay_hours"`
	Symptoms   json.RawMessage `json:"symptoms"`
	Notes      *string         `json:"notes,omitempty"`
	CreatedAt  time.Time       `json:"created_at"`
}

// CreateCheckinRequest holds the data needed to create a symptom check-in.
type CreateCheckinRequest struct {
	DelayHours int              `json:"delay_hours" binding:"required,oneof=6 8 12"`
	Symptoms   []SymptomDetail  `json:"symptoms" binding:"required,dive"`
	Notes      *string          `json:"notes,omitempty" binding:"omitempty,max=1000"`
}

// SymptomEntry represents a standalone symptom entry not linked to a meal.
type SymptomEntry struct {
	ID        uuid.UUID       `json:"id"`
	UserID    uuid.UUID       `json:"user_id"`
	Symptoms  json.RawMessage `json:"symptoms"`
	Notes     *string         `json:"notes,omitempty"`
	EntryTime time.Time       `json:"entry_time"`
	CreatedAt time.Time       `json:"created_at"`
}

// CreateSymptomEntryRequest holds the data needed to create a standalone symptom entry.
type CreateSymptomEntryRequest struct {
	Symptoms  []SymptomDetail `json:"symptoms" binding:"required,dive"`
	Notes     *string         `json:"notes,omitempty" binding:"omitempty,max=1000"`
	EntryTime string          `json:"entry_time" binding:"required"`
}
