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
	ErrSymptomEntryNotFound = errors.New("symptom entry not found")
	ErrCheckinNotFound      = errors.New("check-in not found")
)

type SymptomService struct {
	db  *pgxpool.Pool
	enc *utils.Encryptor
}

func NewSymptomService(db *pgxpool.Pool, enc *utils.Encryptor) *SymptomService {
	return &SymptomService{db: db, enc: enc}
}

func (s *SymptomService) decryptEntry(entry *models.SymptomEntry, encSymptoms string, encNotes *string) {
	decSymptoms, _ := s.enc.Decrypt(encSymptoms)
	if decSymptoms != "" {
		entry.Symptoms = json.RawMessage(decSymptoms)
	} else {
		entry.Symptoms = json.RawMessage("[]")
	}
	if encNotes != nil && *encNotes != "" {
		dec, _ := s.enc.Decrypt(*encNotes)
		entry.Notes = &dec
	}
}

func (s *SymptomService) Create(ctx context.Context, userID uuid.UUID, req models.CreateSymptomEntryRequest) (*models.SymptomEntry, error) {
	entryTime, err := time.Parse(time.RFC3339, req.EntryTime)
	if err != nil {
		return nil, fmt.Errorf("invalid entry_time format (expected RFC3339): %w", err)
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

	var entry models.SymptomEntry
	var dbSymptoms string
	var dbNotes *string
	err = s.db.QueryRow(ctx,
		`INSERT INTO symptom_entries (user_id, symptoms, notes, entry_time)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, user_id, symptoms, notes, entry_time, created_at`,
		userID, encSymptoms, encNotes, entryTime,
	).Scan(&entry.ID, &entry.UserID, &dbSymptoms, &dbNotes, &entry.EntryTime, &entry.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create symptom entry: %w", err)
	}

	s.decryptEntry(&entry, dbSymptoms, dbNotes)
	return &entry, nil
}

func (s *SymptomService) ListByDate(ctx context.Context, userID uuid.UUID, date string) ([]models.SymptomEntry, error) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format (expected YYYY-MM-DD): %w", err)
	}

	rows, err := s.db.Query(ctx,
		`SELECT id, user_id, symptoms, notes, entry_time, created_at
		 FROM symptom_entries WHERE user_id = $1 AND entry_time >= $2 AND entry_time < $3
		 ORDER BY entry_time DESC`,
		userID, parsedDate, parsedDate.Add(24*time.Hour),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list symptom entries: %w", err)
	}
	defer rows.Close()

	var entries []models.SymptomEntry
	for rows.Next() {
		var e models.SymptomEntry
		var encSymptoms string
		var encNotes *string
		if err := rows.Scan(&e.ID, &e.UserID, &encSymptoms, &encNotes, &e.EntryTime, &e.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan symptom entry: %w", err)
		}
		s.decryptEntry(&e, encSymptoms, encNotes)
		entries = append(entries, e)
	}

	if entries == nil {
		entries = []models.SymptomEntry{}
	}
	return entries, nil
}

func (s *SymptomService) Update(ctx context.Context, userID, entryID uuid.UUID, req models.CreateSymptomEntryRequest) (*models.SymptomEntry, error) {
	entryTime, err := time.Parse(time.RFC3339, req.EntryTime)
	if err != nil {
		return nil, fmt.Errorf("invalid entry_time format (expected RFC3339): %w", err)
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

	var entry models.SymptomEntry
	var dbSymptoms string
	var dbNotes *string
	err = s.db.QueryRow(ctx,
		`UPDATE symptom_entries SET symptoms = $1, notes = $2, entry_time = $3
		 WHERE id = $4 AND user_id = $5
		 RETURNING id, user_id, symptoms, notes, entry_time, created_at`,
		encSymptoms, encNotes, entryTime, entryID, userID,
	).Scan(&entry.ID, &entry.UserID, &dbSymptoms, &dbNotes, &entry.EntryTime, &entry.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSymptomEntryNotFound
		}
		return nil, fmt.Errorf("failed to update symptom entry: %w", err)
	}

	s.decryptEntry(&entry, dbSymptoms, dbNotes)
	return &entry, nil
}

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

func (s *SymptomService) UpdateCheckin(ctx context.Context, userID, checkinID uuid.UUID, req models.CreateCheckinRequest) (*models.SymptomCheckin, error) {
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
		`UPDATE symptom_checkins SET delay_hours = $1, symptoms = $2, notes = $3
		 WHERE id = $4 AND meal_id IN (SELECT id FROM meals WHERE user_id = $5)
		 RETURNING id, meal_id, delay_hours, symptoms, notes, created_at`,
		req.DelayHours, encSymptoms, encNotes, checkinID, userID,
	).Scan(&checkin.ID, &checkin.MealID, &checkin.DelayHours, &dbSymptoms, &dbNotes, &checkin.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrCheckinNotFound
		}
		return nil, fmt.Errorf("failed to update check-in: %w", err)
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
