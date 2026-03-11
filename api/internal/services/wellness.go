package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/teyk0o/loupi/api/internal/models"
	"github.com/teyk0o/loupi/api/internal/utils"
)

var (
	ErrWellnessNotFound = errors.New("wellness entry not found")
)

type WellnessService struct {
	db  *pgxpool.Pool
	enc *utils.Encryptor
}

func NewWellnessService(db *pgxpool.Pool, enc *utils.Encryptor) *WellnessService {
	return &WellnessService{db: db, enc: enc}
}

func (s *WellnessService) encryptNotes(notes *string) (*string, error) {
	if notes == nil || *notes == "" {
		return notes, nil
	}
	enc, err := s.enc.Encrypt(*notes)
	if err != nil {
		return nil, err
	}
	return &enc, nil
}

func (s *WellnessService) encryptSport(sport *json.RawMessage) (*string, error) {
	if sport == nil {
		return nil, nil
	}
	enc, err := s.enc.Encrypt(string(*sport))
	if err != nil {
		return nil, err
	}
	return &enc, nil
}

func (s *WellnessService) decryptEntry(entry *models.WellnessEntry, encSport *string, encNotes *string) {
	if encNotes != nil && *encNotes != "" {
		dec, _ := s.enc.Decrypt(*encNotes)
		entry.Notes = &dec
	}
	if encSport != nil && *encSport != "" {
		dec, _ := s.enc.Decrypt(*encSport)
		if dec != "" {
			raw := json.RawMessage(dec)
			entry.Sport = raw
		}
	}
}

func (s *WellnessService) Upsert(ctx context.Context, userID uuid.UUID, req models.CreateWellnessRequest) (*models.WellnessEntry, error) {
	encNotes, err := s.encryptNotes(req.Notes)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt notes: %w", err)
	}

	encSport, err := s.encryptSport(req.Sport)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt sport: %w", err)
	}

	var entry models.WellnessEntry
	var dbSport, dbNotes *string
	err = s.db.QueryRow(ctx,
		`INSERT INTO wellness_entries (user_id, date, stress, mood, energy, sleep_hours, sleep_quality, sport, hydration, notes)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		 ON CONFLICT (user_id, date)
		 DO UPDATE SET stress = $3, mood = $4, energy = $5, sleep_hours = $6, sleep_quality = $7,
		              sport = $8, hydration = $9, notes = $10, updated_at = NOW()
		 RETURNING id, user_id, date, stress, mood, energy, sleep_hours, sleep_quality, sport, hydration, notes, created_at, updated_at`,
		userID, req.Date, req.Stress, req.Mood, req.Energy, req.SleepHours, req.SleepQuality, encSport, req.Hydration, encNotes,
	).Scan(&entry.ID, &entry.UserID, &entry.Date, &entry.Stress, &entry.Mood, &entry.Energy, &entry.SleepHours, &entry.SleepQuality, &dbSport, &entry.Hydration, &dbNotes, &entry.CreatedAt, &entry.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert wellness entry: %w", err)
	}

	s.decryptEntry(&entry, dbSport, dbNotes)
	return &entry, nil
}

func (s *WellnessService) GetByDate(ctx context.Context, userID uuid.UUID, date string) (*models.WellnessEntry, error) {
	var entry models.WellnessEntry
	var dbSport, dbNotes *string
	err := s.db.QueryRow(ctx,
		`SELECT id, user_id, date, stress, mood, energy, sleep_hours, sleep_quality, sport, hydration, notes, created_at, updated_at
		 FROM wellness_entries WHERE user_id = $1 AND date = $2`,
		userID, date,
	).Scan(&entry.ID, &entry.UserID, &entry.Date, &entry.Stress, &entry.Mood, &entry.Energy, &entry.SleepHours, &entry.SleepQuality, &dbSport, &entry.Hydration, &dbNotes, &entry.CreatedAt, &entry.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWellnessNotFound
		}
		return nil, fmt.Errorf("failed to get wellness entry: %w", err)
	}

	s.decryptEntry(&entry, dbSport, dbNotes)
	return &entry, nil
}

func (s *WellnessService) Update(ctx context.Context, userID, entryID uuid.UUID, req models.CreateWellnessRequest) (*models.WellnessEntry, error) {
	encNotes, err := s.encryptNotes(req.Notes)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt notes: %w", err)
	}

	encSport, err := s.encryptSport(req.Sport)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt sport: %w", err)
	}

	var entry models.WellnessEntry
	var dbSport, dbNotes *string
	err = s.db.QueryRow(ctx,
		`UPDATE wellness_entries
		 SET stress = $1, mood = $2, energy = $3, sleep_hours = $4, sleep_quality = $5,
		     sport = $6, hydration = $7, notes = $8, updated_at = NOW()
		 WHERE id = $9 AND user_id = $10
		 RETURNING id, user_id, date, stress, mood, energy, sleep_hours, sleep_quality, sport, hydration, notes, created_at, updated_at`,
		req.Stress, req.Mood, req.Energy, req.SleepHours, req.SleepQuality, encSport, req.Hydration, encNotes, entryID, userID,
	).Scan(&entry.ID, &entry.UserID, &entry.Date, &entry.Stress, &entry.Mood, &entry.Energy, &entry.SleepHours, &entry.SleepQuality, &dbSport, &entry.Hydration, &dbNotes, &entry.CreatedAt, &entry.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrWellnessNotFound
		}
		return nil, fmt.Errorf("failed to update wellness entry: %w", err)
	}

	s.decryptEntry(&entry, dbSport, dbNotes)
	return &entry, nil
}
