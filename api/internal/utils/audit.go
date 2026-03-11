package utils

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AuditLogger writes audit log entries to the database.
type AuditLogger struct {
	db *pgxpool.Pool
}

// NewAuditLogger creates a new AuditLogger.
func NewAuditLogger(db *pgxpool.Pool) *AuditLogger {
	return &AuditLogger{db: db}
}

// Log writes an audit log entry. resourceID can be uuid.Nil if not applicable.
func (a *AuditLogger) Log(ctx context.Context, userID uuid.UUID, action, resourceType string, resourceID uuid.UUID, ipAddress string) error {
	var resID *uuid.UUID
	if resourceID != uuid.Nil {
		resID = &resourceID
	}

	_, err := a.db.Exec(ctx,
		`INSERT INTO audit_logs (user_id, action, resource_type, resource_id, ip_address)
		 VALUES ($1, $2, $3, $4, $5::inet)`,
		userID, action, resourceType, resID, ipAddress,
	)
	if err != nil {
		return fmt.Errorf("failed to write audit log: %w", err)
	}
	return nil
}
