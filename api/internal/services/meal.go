package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/teyk0o/loupi/api/internal/models"
)

// Common errors returned by the meal service.
var (
	ErrMealNotFound = errors.New("meal not found")
)

// MealService handles meal-related business logic.
type MealService struct {
	db *pgxpool.Pool
}

// NewMealService creates a new MealService instance.
func NewMealService(db *pgxpool.Pool) *MealService {
	return &MealService{db: db}
}

// Create adds a new meal for the given user.
func (s *MealService) Create(ctx context.Context, userID uuid.UUID, req models.CreateMealRequest) (*models.Meal, error) {
	mealTime, err := time.Parse(time.RFC3339, req.MealTime)
	if err != nil {
		return nil, fmt.Errorf("invalid meal_time format (expected RFC3339): %w", err)
	}

	var meal models.Meal
	err = s.db.QueryRow(ctx,
		`INSERT INTO meals (user_id, description, category, meal_time)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, user_id, photo_uuid, description, category, meal_time, created_at, updated_at`,
		userID, req.Description, req.Category, mealTime,
	).Scan(&meal.ID, &meal.UserID, &meal.PhotoUUID, &meal.Description, &meal.Category, &meal.MealTime, &meal.CreatedAt, &meal.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create meal: %w", err)
	}

	return &meal, nil
}

// GetByID retrieves a meal by its ID, scoped to the given user.
func (s *MealService) GetByID(ctx context.Context, userID, mealID uuid.UUID) (*models.Meal, error) {
	var meal models.Meal
	err := s.db.QueryRow(ctx,
		`SELECT id, user_id, photo_uuid, description, category, meal_time, created_at, updated_at
		 FROM meals WHERE id = $1 AND user_id = $2`,
		mealID, userID,
	).Scan(&meal.ID, &meal.UserID, &meal.PhotoUUID, &meal.Description, &meal.Category, &meal.MealTime, &meal.CreatedAt, &meal.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrMealNotFound
		}
		return nil, fmt.Errorf("failed to get meal: %w", err)
	}

	return &meal, nil
}

// ListByDate returns all meals for a user on a given date.
func (s *MealService) ListByDate(ctx context.Context, userID uuid.UUID, date string) ([]models.Meal, error) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format (expected YYYY-MM-DD): %w", err)
	}

	startOfDay := parsedDate
	endOfDay := parsedDate.Add(24 * time.Hour)

	rows, err := s.db.Query(ctx,
		`SELECT id, user_id, photo_uuid, description, category, meal_time, created_at, updated_at
		 FROM meals WHERE user_id = $1 AND meal_time >= $2 AND meal_time < $3
		 ORDER BY meal_time DESC`,
		userID, startOfDay, endOfDay,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list meals: %w", err)
	}
	defer rows.Close()

	var meals []models.Meal
	for rows.Next() {
		var meal models.Meal
		if err := rows.Scan(&meal.ID, &meal.UserID, &meal.PhotoUUID, &meal.Description, &meal.Category, &meal.MealTime, &meal.CreatedAt, &meal.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan meal: %w", err)
		}
		meals = append(meals, meal)
	}

	if meals == nil {
		meals = []models.Meal{}
	}

	return meals, nil
}

// Update modifies an existing meal.
func (s *MealService) Update(ctx context.Context, userID, mealID uuid.UUID, req models.UpdateMealRequest) (*models.Meal, error) {
	existing, err := s.GetByID(ctx, userID, mealID)
	if err != nil {
		return nil, err
	}

	desc := existing.Description
	if req.Description != nil {
		desc = *req.Description
	}

	cat := existing.Category
	if req.Category != nil {
		cat = *req.Category
	}

	mealTime := existing.MealTime
	if req.MealTime != nil {
		mealTime, err = time.Parse(time.RFC3339, *req.MealTime)
		if err != nil {
			return nil, fmt.Errorf("invalid meal_time format (expected RFC3339): %w", err)
		}
	}

	var meal models.Meal
	err = s.db.QueryRow(ctx,
		`UPDATE meals SET description = $1, category = $2, meal_time = $3, updated_at = NOW()
		 WHERE id = $4 AND user_id = $5
		 RETURNING id, user_id, photo_uuid, description, category, meal_time, created_at, updated_at`,
		desc, cat, mealTime, mealID, userID,
	).Scan(&meal.ID, &meal.UserID, &meal.PhotoUUID, &meal.Description, &meal.Category, &meal.MealTime, &meal.CreatedAt, &meal.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update meal: %w", err)
	}

	return &meal, nil
}

// Delete removes a meal by its ID, scoped to the given user.
func (s *MealService) Delete(ctx context.Context, userID, mealID uuid.UUID) error {
	result, err := s.db.Exec(ctx, "DELETE FROM meals WHERE id = $1 AND user_id = $2", mealID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete meal: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrMealNotFound
	}
	return nil
}

// GetCheckins returns all symptom check-ins for a meal.
func (s *MealService) GetCheckins(ctx context.Context, userID, mealID uuid.UUID) ([]models.SymptomCheckin, error) {
	// Verify the meal belongs to the user
	_, err := s.GetByID(ctx, userID, mealID)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(ctx,
		`SELECT id, meal_id, delay_hours, symptoms, notes, created_at
		 FROM symptom_checkins WHERE meal_id = $1
		 ORDER BY delay_hours ASC`,
		mealID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list check-ins: %w", err)
	}
	defer rows.Close()

	var checkins []models.SymptomCheckin
	for rows.Next() {
		var c models.SymptomCheckin
		if err := rows.Scan(&c.ID, &c.MealID, &c.DelayHours, &c.Symptoms, &c.Notes, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan check-in: %w", err)
		}
		checkins = append(checkins, c)
	}

	if checkins == nil {
		checkins = []models.SymptomCheckin{}
	}

	return checkins, nil
}

// CreateCheckin adds a symptom check-in to a meal.
func (s *MealService) CreateCheckin(ctx context.Context, userID, mealID uuid.UUID, req models.CreateCheckinRequest) (*models.SymptomCheckin, error) {
	// Verify the meal belongs to the user
	_, err := s.GetByID(ctx, userID, mealID)
	if err != nil {
		return nil, err
	}

	symptomsJSON, err := json.Marshal(req.Symptoms)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal symptoms: %w", err)
	}

	var checkin models.SymptomCheckin
	err = s.db.QueryRow(ctx,
		`INSERT INTO symptom_checkins (meal_id, delay_hours, symptoms, notes)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, meal_id, delay_hours, symptoms, notes, created_at`,
		mealID, req.DelayHours, symptomsJSON, req.Notes,
	).Scan(&checkin.ID, &checkin.MealID, &checkin.DelayHours, &checkin.Symptoms, &checkin.Notes, &checkin.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create check-in: %w", err)
	}

	return &checkin, nil
}
