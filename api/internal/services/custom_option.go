package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/teyk0o/loupi/api/internal/models"
)

// Common errors returned by the custom option service.
var (
	ErrOptionNotFound      = errors.New("custom option not found")
	ErrOptionAlreadyExists = errors.New("option value already exists for this category")
	ErrInvalidCategory     = errors.New("invalid option category")
)

// validCategories defines the allowed option categories.
var validCategories = map[string]bool{
	"symptom_type":  true,
	"meal_category": true,
	"sport_type":    true,
}

// CustomOptionService handles custom option business logic.
type CustomOptionService struct {
	db *pgxpool.Pool
}

// NewCustomOptionService creates a new CustomOptionService instance.
func NewCustomOptionService(db *pgxpool.Pool) *CustomOptionService {
	return &CustomOptionService{db: db}
}

// ListByCategory retrieves all custom options for a user in a given category, ordered by sort_order.
func (s *CustomOptionService) ListByCategory(ctx context.Context, userID uuid.UUID, category string) ([]models.CustomOption, error) {
	if !validCategories[category] {
		return nil, ErrInvalidCategory
	}

	rows, err := s.db.Query(ctx,
		`SELECT id, user_id, category, value, label, emoji, sort_order, created_at
		 FROM user_custom_options
		 WHERE user_id = $1 AND category = $2
		 ORDER BY sort_order, created_at`,
		userID, category,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list custom options: %w", err)
	}
	defer rows.Close()

	var options []models.CustomOption
	for rows.Next() {
		var opt models.CustomOption
		if err := rows.Scan(&opt.ID, &opt.UserID, &opt.Category, &opt.Value, &opt.Label, &opt.Emoji, &opt.SortOrder, &opt.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan custom option: %w", err)
		}
		options = append(options, opt)
	}

	return options, nil
}

// Create adds a new custom option for a user in the given category.
func (s *CustomOptionService) Create(ctx context.Context, userID uuid.UUID, category string, req models.CreateCustomOptionRequest) (*models.CustomOption, error) {
	if !validCategories[category] {
		return nil, ErrInvalidCategory
	}

	// Normalize value to snake_case lowercase.
	normalizedValue := strings.ToLower(strings.TrimSpace(req.Value))
	normalizedValue = strings.ReplaceAll(normalizedValue, " ", "_")

	// Get the next sort_order value.
	var maxOrder int
	err := s.db.QueryRow(ctx,
		`SELECT COALESCE(MAX(sort_order), 0) FROM user_custom_options WHERE user_id = $1 AND category = $2`,
		userID, category,
	).Scan(&maxOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to get max sort order: %w", err)
	}

	var opt models.CustomOption
	err = s.db.QueryRow(ctx,
		`INSERT INTO user_custom_options (user_id, category, value, label, emoji, sort_order)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 RETURNING id, user_id, category, value, label, emoji, sort_order, created_at`,
		userID, category, normalizedValue, strings.TrimSpace(req.Label), req.Emoji, maxOrder+1,
	).Scan(&opt.ID, &opt.UserID, &opt.Category, &opt.Value, &opt.Label, &opt.Emoji, &opt.SortOrder, &opt.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrOptionAlreadyExists
		}
		return nil, fmt.Errorf("failed to create custom option: %w", err)
	}

	return &opt, nil
}

// Update modifies an existing custom option.
func (s *CustomOptionService) Update(ctx context.Context, userID, optionID uuid.UUID, req models.UpdateCustomOptionRequest) (*models.CustomOption, error) {
	var opt models.CustomOption
	err := s.db.QueryRow(ctx,
		`UPDATE user_custom_options
		 SET label = COALESCE($1, label),
		     emoji = COALESCE($2, emoji),
		     sort_order = COALESCE($3, sort_order)
		 WHERE id = $4 AND user_id = $5
		 RETURNING id, user_id, category, value, label, emoji, sort_order, created_at`,
		req.Label, req.Emoji, req.SortOrder, optionID, userID,
	).Scan(&opt.ID, &opt.UserID, &opt.Category, &opt.Value, &opt.Label, &opt.Emoji, &opt.SortOrder, &opt.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOptionNotFound
		}
		return nil, fmt.Errorf("failed to update custom option: %w", err)
	}

	return &opt, nil
}

// Delete removes a custom option by ID.
func (s *CustomOptionService) Delete(ctx context.Context, userID, optionID uuid.UUID) error {
	result, err := s.db.Exec(ctx,
		`DELETE FROM user_custom_options WHERE id = $1 AND user_id = $2`,
		optionID, userID,
	)
	if err != nil {
		return fmt.Errorf("failed to delete custom option: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrOptionNotFound
	}
	return nil
}

// Reorder updates the sort_order for all options in a category based on the provided ID order.
func (s *CustomOptionService) Reorder(ctx context.Context, userID uuid.UUID, category string, req models.ReorderCustomOptionsRequest) error {
	if !validCategories[category] {
		return ErrInvalidCategory
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	for i, id := range req.OrderedIDs {
		result, err := tx.Exec(ctx,
			`UPDATE user_custom_options SET sort_order = $1 WHERE id = $2 AND user_id = $3 AND category = $4`,
			i+1, id, userID, category,
		)
		if err != nil {
			return fmt.Errorf("failed to update sort order: %w", err)
		}
		if result.RowsAffected() == 0 {
			return ErrOptionNotFound
		}
	}

	return tx.Commit(ctx)
}
