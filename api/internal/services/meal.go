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
	"github.com/teyk0o/loupi/api/internal/utils"
)

var (
	ErrMealNotFound = errors.New("meal not found")
)

type MealService struct {
	db  *pgxpool.Pool
	enc *utils.Encryptor
}

func NewMealService(db *pgxpool.Pool, enc *utils.Encryptor) *MealService {
	return &MealService{db: db, enc: enc}
}

func (s *MealService) Create(ctx context.Context, userID uuid.UUID, req models.CreateMealRequest) (*models.Meal, error) {
	mealTime, err := time.Parse(time.RFC3339, req.MealTime)
	if err != nil {
		return nil, fmt.Errorf("invalid meal_time format (expected RFC3339): %w", err)
	}

	encDesc, err := s.enc.Encrypt(req.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt description: %w", err)
	}

	var meal models.Meal
	err = s.db.QueryRow(ctx,
		`INSERT INTO meals (user_id, description, category, meal_time)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, user_id, photo_uuid, description, category, meal_time, created_at, updated_at`,
		userID, encDesc, req.Category, mealTime,
	).Scan(&meal.ID, &meal.UserID, &meal.PhotoUUID, &meal.Description, &meal.Category, &meal.MealTime, &meal.CreatedAt, &meal.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create meal: %w", err)
	}

	meal.Description, _ = s.enc.Decrypt(meal.Description)
	return &meal, nil
}

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

	meal.Description, _ = s.enc.Decrypt(meal.Description)
	return &meal, nil
}

func (s *MealService) ListByDate(ctx context.Context, userID uuid.UUID, date string) ([]models.Meal, error) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format (expected YYYY-MM-DD): %w", err)
	}

	rows, err := s.db.Query(ctx,
		`SELECT id, user_id, photo_uuid, description, category, meal_time, created_at, updated_at
		 FROM meals WHERE user_id = $1 AND meal_time >= $2 AND meal_time < $3
		 ORDER BY meal_time DESC`,
		userID, parsedDate, parsedDate.Add(24*time.Hour),
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
		meal.Description, _ = s.enc.Decrypt(meal.Description)
		meals = append(meals, meal)
	}

	if meals == nil {
		meals = []models.Meal{}
	}
	return meals, nil
}

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

	encDesc, err := s.enc.Encrypt(desc)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt description: %w", err)
	}

	var meal models.Meal
	err = s.db.QueryRow(ctx,
		`UPDATE meals SET description = $1, category = $2, meal_time = $3, updated_at = NOW()
		 WHERE id = $4 AND user_id = $5
		 RETURNING id, user_id, photo_uuid, description, category, meal_time, created_at, updated_at`,
		encDesc, cat, mealTime, mealID, userID,
	).Scan(&meal.ID, &meal.UserID, &meal.PhotoUUID, &meal.Description, &meal.Category, &meal.MealTime, &meal.CreatedAt, &meal.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update meal: %w", err)
	}

	meal.Description, _ = s.enc.Decrypt(meal.Description)
	return &meal, nil
}

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

func (s *MealService) GetCheckins(ctx context.Context, userID, mealID uuid.UUID) ([]models.SymptomCheckin, error) {
	if _, err := s.GetByID(ctx, userID, mealID); err != nil {
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
		var encSymptoms string
		var encNotes *string
		if err := rows.Scan(&c.ID, &c.MealID, &c.DelayHours, &encSymptoms, &encNotes, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan check-in: %w", err)
		}
		decSymptoms, _ := s.enc.Decrypt(encSymptoms)
		if decSymptoms != "" {
			c.Symptoms = json.RawMessage(decSymptoms)
		} else {
			c.Symptoms = json.RawMessage("[]")
		}
		if encNotes != nil && *encNotes != "" {
			dec, _ := s.enc.Decrypt(*encNotes)
			c.Notes = &dec
		}
		checkins = append(checkins, c)
	}

	if checkins == nil {
		checkins = []models.SymptomCheckin{}
	}
	return checkins, nil
}

func (s *MealService) CreateCheckin(ctx context.Context, userID, mealID uuid.UUID, req models.CreateCheckinRequest) (*models.SymptomCheckin, error) {
	if _, err := s.GetByID(ctx, userID, mealID); err != nil {
		return nil, err
	}

	encSymptoms, err := s.enc.EncryptJSON(req.Symptoms)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt symptoms: %w", err)
	}

	var encNotes *string
	if req.Notes != nil && *req.Notes != "" {
		enc, err := s.enc.Encrypt(*req.Notes)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt notes: %w", err)
		}
		encNotes = &enc
	}

	var checkin models.SymptomCheckin
	var dbSymptoms string
	var dbNotes *string
	err = s.db.QueryRow(ctx,
		`INSERT INTO symptom_checkins (meal_id, delay_hours, symptoms, notes)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, meal_id, delay_hours, symptoms, notes, created_at`,
		mealID, req.DelayHours, encSymptoms, encNotes,
	).Scan(&checkin.ID, &checkin.MealID, &checkin.DelayHours, &dbSymptoms, &dbNotes, &checkin.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create check-in: %w", err)
	}

	decSymptoms, _ := s.enc.Decrypt(dbSymptoms)
	if decSymptoms != "" {
		checkin.Symptoms = json.RawMessage(decSymptoms)
	} else {
		checkin.Symptoms = json.RawMessage("[]")
	}
	if dbNotes != nil && *dbNotes != "" {
		dec, _ := s.enc.Decrypt(*dbNotes)
		checkin.Notes = &dec
	}

	return &checkin, nil
}
