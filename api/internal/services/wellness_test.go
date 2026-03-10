package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/teyk0o/loupi/api/internal/models"
)

// setupTestWellnessEnv creates a test user and returns the wellness service.
func setupTestWellnessEnv(t *testing.T) (*WellnessService, uuid.UUID, func()) {
	t.Helper()
	pool := setupTestDB(t)
	rdb := setupTestRedis(t)
	cfg := testConfig()
	enc := testEncryptor(t)
	authSvc := NewAuthService(pool, rdb, cfg)
	wellnessSvc := NewWellnessService(pool, enc)
	ctx := context.Background()

	email := "test-wellness-" + uuid.New().String()[:8] + "@loupi.test"
	resp, err := authSvc.Register(ctx, models.RegisterRequest{Email: email, Password: "SecurePass123#"})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	cleanup := func() {
		_ = authSvc.DeleteAccount(ctx, resp.User.ID)
		rdb.Close()
		pool.Close()
	}

	return wellnessSvc, resp.User.ID, cleanup
}

func TestWellness_Upsert_Create(t *testing.T) {
	svc, userID, cleanup := setupTestWellnessEnv(t)
	defer cleanup()

	stress := 3
	mood := 4
	hydration := 8
	entry, err := svc.Upsert(context.Background(), userID, models.CreateWellnessRequest{
		Date:      time.Now().Format("2006-01-02"),
		Stress:    &stress,
		Mood:      &mood,
		Hydration: &hydration,
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if *entry.Stress != 3 {
		t.Errorf("expected stress 3, got %d", *entry.Stress)
	}
	if *entry.Mood != 4 {
		t.Errorf("expected mood 4, got %d", *entry.Mood)
	}
	if *entry.Hydration != 8 {
		t.Errorf("expected hydration 8, got %d", *entry.Hydration)
	}
}

func TestWellness_Upsert_Update(t *testing.T) {
	svc, userID, cleanup := setupTestWellnessEnv(t)
	defer cleanup()

	date := time.Now().Format("2006-01-02")
	stress1 := 2
	stress2 := 5

	// Create
	entry1, err := svc.Upsert(context.Background(), userID, models.CreateWellnessRequest{
		Date:   date,
		Stress: &stress1,
	})
	if err != nil {
		t.Fatalf("first upsert failed: %v", err)
	}

	// Update (same date)
	entry2, err := svc.Upsert(context.Background(), userID, models.CreateWellnessRequest{
		Date:   date,
		Stress: &stress2,
	})
	if err != nil {
		t.Fatalf("second upsert failed: %v", err)
	}

	// Should be the same entry (upsert)
	if entry1.ID != entry2.ID {
		t.Errorf("expected same ID after upsert, got %s and %s", entry1.ID, entry2.ID)
	}
	if *entry2.Stress != 5 {
		t.Errorf("expected stress 5 after update, got %d", *entry2.Stress)
	}
}

func TestWellness_GetByDate(t *testing.T) {
	svc, userID, cleanup := setupTestWellnessEnv(t)
	defer cleanup()

	date := time.Now().Format("2006-01-02")
	energy := 4
	_, err := svc.Upsert(context.Background(), userID, models.CreateWellnessRequest{
		Date:   date,
		Energy: &energy,
	})
	if err != nil {
		t.Fatalf("upsert failed: %v", err)
	}

	entry, err := svc.GetByDate(context.Background(), userID, date)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if *entry.Energy != 4 {
		t.Errorf("expected energy 4, got %d", *entry.Energy)
	}
}

func TestWellness_GetByDate_NotFound(t *testing.T) {
	svc, userID, cleanup := setupTestWellnessEnv(t)
	defer cleanup()

	_, err := svc.GetByDate(context.Background(), userID, "2020-01-01")
	if err != ErrWellnessNotFound {
		t.Fatalf("expected ErrWellnessNotFound, got: %v", err)
	}
}

func TestWellness_Update(t *testing.T) {
	svc, userID, cleanup := setupTestWellnessEnv(t)
	defer cleanup()

	date := time.Now().Format("2006-01-02")
	mood := 3
	created, err := svc.Upsert(context.Background(), userID, models.CreateWellnessRequest{
		Date: date,
		Mood: &mood,
	})
	if err != nil {
		t.Fatalf("upsert failed: %v", err)
	}

	newMood := 5
	sleepHours := float32(7.5)
	updated, err := svc.Update(context.Background(), userID, created.ID, models.CreateWellnessRequest{
		Date:       date,
		Mood:       &newMood,
		SleepHours: &sleepHours,
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if *updated.Mood != 5 {
		t.Errorf("expected mood 5, got %d", *updated.Mood)
	}
	if *updated.SleepHours != 7.5 {
		t.Errorf("expected sleep_hours 7.5, got %f", *updated.SleepHours)
	}
}

func TestWellness_Update_NotFound(t *testing.T) {
	svc, userID, cleanup := setupTestWellnessEnv(t)
	defer cleanup()

	mood := 3
	_, err := svc.Update(context.Background(), userID, uuid.New(), models.CreateWellnessRequest{
		Date: time.Now().Format("2006-01-02"),
		Mood: &mood,
	})
	if err != ErrWellnessNotFound {
		t.Fatalf("expected ErrWellnessNotFound, got: %v", err)
	}
}
