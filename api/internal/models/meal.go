package models

import (
	"time"

	"github.com/google/uuid"
)

// Meal represents a logged meal entry.
type Meal struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	PhotoUUID   *uuid.UUID `json:"photo_uuid,omitempty"`
	Description string     `json:"description"`
	Category    string     `json:"category"`
	MealTime    time.Time  `json:"meal_time"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// CreateMealRequest holds the data needed to create a new meal.
type CreateMealRequest struct {
	Description string `json:"description" binding:"required,max=1000"`
	Category    string `json:"category" binding:"required,oneof=homemade restaurant takeout snack fast_food cafeteria family friends other"`
	MealTime    string `json:"meal_time" binding:"required"`
}

// UpdateMealRequest holds the data needed to update an existing meal.
type UpdateMealRequest struct {
	Description *string `json:"description,omitempty" binding:"omitempty,max=1000"`
	Category    *string `json:"category,omitempty" binding:"omitempty,oneof=homemade restaurant takeout snack fast_food cafeteria family friends other"`
	MealTime    *string `json:"meal_time,omitempty"`
}
