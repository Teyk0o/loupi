package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/teyk0o/loupi/api/internal/models"
)

// Common errors returned by the wellness service.
var (
	ErrWellnessNotFound = errors.New("wellness entry not found")
)

// WellnessService handles wellness entry business logic.
type WellnessService struct {
	db *pgxpool.Pool
}

// NewWellnessService creates a new WellnessService instance.
func NewWellnessService(db *pgxpool.Pool) *WellnessService {
	return &WellnessService{db: db}
}

// Upsert creates or updates a wellness entry for the given user and date.
// Uses PostgreSQL ON CONFLICT to handle the unique (user_id, date) constraint.
func (s *WellnessService) Upsert(ctx context.Context, userID uuid.UUID, req models.CreateWellnessRequest) (*models.WellnessEntry, error) {
	var entry models.WellnessEntry
	err := s.db.QueryRow(ctx,
		`INSERT INTO wellness_entries (user_id, date, stress, mood, energy, sleep_hours, sleep_quality, sport, hydration, notes)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 ON CONFLICT (user_id, date)
		 DO UPDATE SET stress = $3, mood = $4, energy = $5, sleep_hours = $6, sleep_quality = $7,
		              sport = $8, hydration = $9, notes = $10, updated_at = NOW()
		 RETURNING id, user_id, date, stress, mood, energy, sleep_hours, sleep_quality, sport, hydration, notes, created_at, updated_at`,
		userID, req.Date, req.Stress, req.Mood, req.Energy, req.SleepHours, req.SleepQuality, req.Sport, req.Hydration, req.Notes,
	).Scan(&entry.ID, &entry.UserID, &entry.Date, &entry.Stress, &entry.Mood, &entry.Energy, &entry.SleepHours, &entry.SleepQuality, &entry.Sport, &entry.Hydration, &entry.Notes, &entry.CreatedAt, &entry.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert wellness entry: %w", err)
	}

	return &entry, nil
}

// GetByDate retrieves a wellness entry for a user on a given date.
func (s *WellnessService) GetByDate(ctx context.Context, userID uuid.UUID, date string) (*models.WellnessEntry, error) {
	var entry models.WellnessEntry
	err := s.db.QueryRow(ctx,
		`SELECT id, user_id, date, stress, mood, energy, sleep_hours, sleep_quality, sport, hydration, notes, created_at, updated_at
		 FROM wellness_entries WHERE user_id = $1 AND date = $2`,
		userID, date,
	).Scan(&entry.ID, &entry.UserID, &entry.Date, &entry.Stress, &entry.Mood, &entry.Energy, &entry.SleepHours, &entry.SleepQuality, &entry.Sport, &entry.Hydration, &entry.Notes, &entry.CreatedAt, &entry.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWellnessNotFound
		}
		return nil, fmt.Errorf("failed to get wellness entry: %w", err)
	}

	return &entry, nil
}

// Update modifies an existing wellness entry by ID.
func (s *WellnessService) Update(ctx context.Context, userID, entryID uuid.UUID, req models.CreateWellnessRequest) (*models.WellnessEntry, error) {
	var entry models.WellnessEntry
	err := s.db.QueryRow(ctx,
		`UPDATE wellness_entries
		 SET stress = $1, mood = $2, energy = $3, sleep_hours = $4, sleep_quality = $5,
		     sport = $6, hydration = $7, notes = $8, updated_at = NOW()
		 WHERE id = $9 AND user_id = $10
		 RETURNING id, user_id, date, stress, mood, energy, sleep_hours, sleep_quality, sport, hydration, notes, created_at, updated_at`,
		req.Stress, req.Mood, req.Energy, req.SleepHours, req.SleepQuality, req.Sport, req.Hydration, req.Notes, entryID, userID,
	).Scan(&entry.ID, &entry.UserID, &entry.Date, &entry.Stress, &entry.Mood, &entry.Energy, &entry.SleepHours, &entry.SleepQuality, &entry.Sport, &entry.Hydration, &entry.Notes, &entry.CreatedAt, &entry.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWellnessNotFound
		}
		return nil, fmt.Errorf("failed to update wellness entry: %w", err)
	}

	return &entry, nil
}
