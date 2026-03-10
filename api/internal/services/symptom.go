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

// Common errors returned by the symptom service.
var (
	ErrSymptomEntryNotFound = errors.New("symptom entry not found")
	ErrCheckinNotFound      = errors.New("check-in not found")
)

// SymptomService handles standalone symptom entry business logic.
type SymptomService struct {
	db *pgxpool.Pool
}

// NewSymptomService creates a new SymptomService instance.
func NewSymptomService(db *pgxpool.Pool) *SymptomService {
	return &SymptomService{db: db}
}

// Create adds a new standalone symptom entry for the given user.
func (s *SymptomService) Create(ctx context.Context, userID uuid.UUID, req models.CreateSymptomEntryRequest) (*models.SymptomEntry, error) {
	entryTime, err := time.Parse(time.RFC3339, req.EntryTime)
	if err != nil {
		return nil, fmt.Errorf("invalid entry_time format (expected RFC3339): %w", err)
	}

	symptomsJSON, err := json.Marshal(req.Symptoms)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal symptoms: %w", err)
	}

	var entry models.SymptomEntry
	err = s.db.QueryRow(ctx,
		`INSERT INTO symptom_entries (user_id, symptoms, notes, entry_time)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, user_id, symptoms, notes, entry_time, created_at`,
		userID, symptomsJSON, req.Notes, entryTime,
	).Scan(&entry.ID, &entry.UserID, &entry.Symptoms, &entry.Notes, &entry.EntryTime, &entry.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create symptom entry: %w", err)
	}

	return &entry, nil
}

// ListByDate returns all standalone symptom entries for a user on a given date.
func (s *SymptomService) ListByDate(ctx context.Context, userID uuid.UUID, date string) ([]models.SymptomEntry, error) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format (expected YYYY-MM-DD): %w", err)
	}

	startOfDay := parsedDate
	endOfDay := parsedDate.Add(24 * time.Hour)

	rows, err := s.db.Query(ctx,
		`SELECT id, user_id, symptoms, notes, entry_time, created_at
		 FROM symptom_entries WHERE user_id = $1 AND entry_time >= $2 AND entry_time < $3
		 ORDER BY entry_time DESC`,
		userID, startOfDay, endOfDay,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list symptom entries: %w", err)
	}
	defer rows.Close()

	var entries []models.SymptomEntry
	for rows.Next() {
		var e models.SymptomEntry
		if err := rows.Scan(&e.ID, &e.UserID, &e.Symptoms, &e.Notes, &e.EntryTime, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan symptom entry: %w", err)
		}
		entries = append(entries, e)
	}

	if entries == nil {
		entries = []models.SymptomEntry{}
	}

	return entries, nil
}

// Update modifies an existing standalone symptom entry.
func (s *SymptomService) Update(ctx context.Context, userID, entryID uuid.UUID, req models.CreateSymptomEntryRequest) (*models.SymptomEntry, error) {
	entryTime, err := time.Parse(time.RFC3339, req.EntryTime)
	if err != nil {
		return nil, fmt.Errorf("invalid entry_time format (expected RFC3339): %w", err)
	}

	symptomsJSON, err := json.Marshal(req.Symptoms)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal symptoms: %w", err)
	}

	var entry models.SymptomEntry
	err = s.db.QueryRow(ctx,
		`UPDATE symptom_entries SET symptoms = $1, notes = $2, entry_time = $3
		 WHERE id = $4 AND user_id = $5
		 RETURNING id, user_id, symptoms, notes, entry_time, created_at`,
		symptomsJSON, req.Notes, entryTime, entryID, userID,
	).Scan(&entry.ID, &entry.UserID, &entry.Symptoms, &entry.Notes, &entry.EntryTime, &entry.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSymptomEntryNotFound
		}
		return nil, fmt.Errorf("failed to update symptom entry: %w", err)
	}

	return &entry, nil
}

// Delete removes a standalone symptom entry.
func (s *SymptomService) Delete(ctx context.Context, userID, entryID uuid.UUID) error {
	result, err := s.db.Exec(ctx, "DELETE FROM symptom_entries WHERE id = $1 AND user_id = $2", entryID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete symptom entry: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrSymptomEntryNotFound
	}
	return nil
}

// UpdateCheckin modifies an existing symptom check-in.
func (s *SymptomService) UpdateCheckin(ctx context.Context, userID, checkinID uuid.UUID, req models.CreateCheckinRequest) (*models.SymptomCheckin, error) {
	symptomsJSON, err := json.Marshal(req.Symptoms)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal symptoms: %w", err)
	}

	// Verify ownership via the meal's user_id
	var checkin models.SymptomCheckin
	err = s.db.QueryRow(ctx,
		`UPDATE symptom_checkins SET delay_hours = $1, symptoms = $2, notes = $3
		 WHERE id = $4 AND meal_id IN (SELECT id FROM meals WHERE user_id = $5)
		 RETURNING id, meal_id, delay_hours, symptoms, notes, created_at`,
		req.DelayHours, symptomsJSON, req.Notes, checkinID, userID,
	).Scan(&checkin.ID, &checkin.MealID, &checkin.DelayHours, &checkin.Symptoms, &checkin.Notes, &checkin.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCheckinNotFound
		}
		return nil, fmt.Errorf("failed to update check-in: %w", err)
	}

	return &checkin, nil
}

// DeleteCheckin removes a symptom check-in.
func (s *SymptomService) DeleteCheckin(ctx context.Context, userID, checkinID uuid.UUID) error {
	result, err := s.db.Exec(ctx,
		`DELETE FROM symptom_checkins WHERE id = $1 AND meal_id IN (SELECT id FROM meals WHERE user_id = $2)`,
		checkinID, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete check-in: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrCheckinNotFound
	}
	return nil
}
