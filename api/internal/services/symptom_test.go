package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/teyk0o/loupi/api/internal/models"
)

// setupTestSymptomEnv creates a test user and returns services and cleanup func.
func setupTestSymptomEnv(t *testing.T) (*SymptomService, *MealService, uuid.UUID, func()) {
	t.Helper()
	pool := setupTestDB(t)
	rdb := setupTestRedis(t)
	cfg := testConfig()
	enc := testEncryptor(t)
	authSvc := NewAuthService(pool, rdb, cfg)
	symptomSvc := NewSymptomService(pool, enc)
	mealSvc := NewMealService(pool, enc)
	ctx := context.Background()

	email := "test-symptom-" + uuid.New().String()[:8] + "@loupi.test"
	resp, err := authSvc.Register(ctx, models.RegisterRequest{Email: email, Password: "SecurePass123#"})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	cleanup := func() {
		_ = authSvc.DeleteAccount(ctx, resp.User.ID)
		rdb.Close()
		pool.Close()
	}

	return symptomSvc, mealSvc, resp.User.ID, cleanup
}

func TestSymptomEntry_Create(t *testing.T) {
	svc, _, userID, cleanup := setupTestSymptomEnv(t)
	defer cleanup()

	entry, err := svc.Create(context.Background(), userID, models.CreateSymptomEntryRequest{
		Symptoms:  []models.SymptomDetail{{Type: "bloating", Severity: 3}, {Type: "gas", Severity: 2}},
		EntryTime: time.Now().Format(time.RFC3339),
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if entry.UserID != userID {
		t.Errorf("expected user ID %s, got %s", userID, entry.UserID)
	}
}

func TestSymptomEntry_Create_InvalidTime(t *testing.T) {
	svc, _, userID, cleanup := setupTestSymptomEnv(t)
	defer cleanup()

	_, err := svc.Create(context.Background(), userID, models.CreateSymptomEntryRequest{
		Symptoms:  []models.SymptomDetail{{Type: "nausea", Severity: 1}},
		EntryTime: "invalid",
	})
	if err == nil {
		t.Fatal("expected error for invalid entry_time")
	}
}

func TestSymptomEntry_ListByDate(t *testing.T) {
	svc, _, userID, cleanup := setupTestSymptomEnv(t)
	defer cleanup()

	now := time.Now()
	_, _ = svc.Create(context.Background(), userID, models.CreateSymptomEntryRequest{
		Symptoms:  []models.SymptomDetail{{Type: "fatigue", Severity: 4}},
		EntryTime: now.Format(time.RFC3339),
	})

	entries, err := svc.ListByDate(context.Background(), userID, now.Format("2006-01-02"))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(entries) < 1 {
		t.Errorf("expected at least 1 entry, got %d", len(entries))
	}
}

func TestSymptomEntry_ListByDate_Empty(t *testing.T) {
	svc, _, userID, cleanup := setupTestSymptomEnv(t)
	defer cleanup()

	entries, err := svc.ListByDate(context.Background(), userID, "2020-01-01")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestSymptomEntry_Update(t *testing.T) {
	svc, _, userID, cleanup := setupTestSymptomEnv(t)
	defer cleanup()

	created, err := svc.Create(context.Background(), userID, models.CreateSymptomEntryRequest{
		Symptoms:  []models.SymptomDetail{{Type: "nausea", Severity: 2}},
		EntryTime: time.Now().Format(time.RFC3339),
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	updated, err := svc.Update(context.Background(), userID, created.ID, models.CreateSymptomEntryRequest{
		Symptoms:  []models.SymptomDetail{{Type: "nausea", Severity: 5}, {Type: "cramps", Severity: 3}},
		EntryTime: time.Now().Format(time.RFC3339),
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if updated.ID != created.ID {
		t.Errorf("expected same ID %s, got %s", created.ID, updated.ID)
	}
}

func TestSymptomEntry_Update_NotFound(t *testing.T) {
	svc, _, userID, cleanup := setupTestSymptomEnv(t)
	defer cleanup()

	_, err := svc.Update(context.Background(), userID, uuid.New(), models.CreateSymptomEntryRequest{
		Symptoms:  []models.SymptomDetail{{Type: "nausea", Severity: 1}},
		EntryTime: time.Now().Format(time.RFC3339),
	})
	if err != ErrSymptomEntryNotFound {
		t.Fatalf("expected ErrSymptomEntryNotFound, got: %v", err)
	}
}

func TestSymptomEntry_Delete(t *testing.T) {
	svc, _, userID, cleanup := setupTestSymptomEnv(t)
	defer cleanup()

	created, err := svc.Create(context.Background(), userID, models.CreateSymptomEntryRequest{
		Symptoms:  []models.SymptomDetail{{Type: "heartburn", Severity: 1}},
		EntryTime: time.Now().Format(time.RFC3339),
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	err = svc.Delete(context.Background(), userID, created.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestSymptomEntry_Delete_NotFound(t *testing.T) {
	svc, _, userID, cleanup := setupTestSymptomEnv(t)
	defer cleanup()

	err := svc.Delete(context.Background(), userID, uuid.New())
	if err != ErrSymptomEntryNotFound {
		t.Fatalf("expected ErrSymptomEntryNotFound, got: %v", err)
	}
}

func TestCheckin_UpdateAndDelete(t *testing.T) {
	symptomSvc, mealSvc, userID, cleanup := setupTestSymptomEnv(t)
	defer cleanup()
	ctx := context.Background()

	// Create a meal first
	meal, err := mealSvc.Create(ctx, userID, models.CreateMealRequest{
		Description: "Test meal for checkin",
		Category:    "homemade",
		MealTime:    time.Now().Format(time.RFC3339),
	})
	if err != nil {
		t.Fatalf("create meal failed: %v", err)
	}

	// Create a checkin
	checkin, err := mealSvc.CreateCheckin(ctx, userID, meal.ID, models.CreateCheckinRequest{
		DelayHours: 6,
		Symptoms:   []models.SymptomDetail{{Type: "bloating", Severity: 2}},
	})
	if err != nil {
		t.Fatalf("create checkin failed: %v", err)
	}

	// Update the checkin
	updated, err := symptomSvc.UpdateCheckin(ctx, userID, checkin.ID, models.CreateCheckinRequest{
		DelayHours: 8,
		Symptoms:   []models.SymptomDetail{{Type: "bloating", Severity: 4}, {Type: "gas", Severity: 3}},
	})
	if err != nil {
		t.Fatalf("update checkin failed: %v", err)
	}
	if updated.DelayHours != 8 {
		t.Errorf("expected delay_hours 8, got %d", updated.DelayHours)
	}

	// Delete the checkin
	err = symptomSvc.DeleteCheckin(ctx, userID, checkin.ID)
	if err != nil {
		t.Fatalf("delete checkin failed: %v", err)
	}

	// Verify deleted
	err = symptomSvc.DeleteCheckin(ctx, userID, checkin.ID)
	if err != ErrCheckinNotFound {
		t.Fatalf("expected ErrCheckinNotFound, got: %v", err)
	}
}

func TestCheckin_Update_WrongUser(t *testing.T) {
	symptomSvc, mealSvc, userID, cleanup := setupTestSymptomEnv(t)
	defer cleanup()
	ctx := context.Background()

	meal, err := mealSvc.Create(ctx, userID, models.CreateMealRequest{
		Description: "Meal",
		Category:    "homemade",
		MealTime:    time.Now().Format(time.RFC3339),
	})
	if err != nil {
		t.Fatalf("create meal failed: %v", err)
	}

	checkin, err := mealSvc.CreateCheckin(ctx, userID, meal.ID, models.CreateCheckinRequest{
		DelayHours: 6,
		Symptoms:   []models.SymptomDetail{{Type: "nausea", Severity: 1}},
	})
	if err != nil {
		t.Fatalf("create checkin failed: %v", err)
	}

	// Another user should not be able to update
	_, err = symptomSvc.UpdateCheckin(ctx, uuid.New(), checkin.ID, models.CreateCheckinRequest{
		DelayHours: 12,
		Symptoms:   []models.SymptomDetail{{Type: "nausea", Severity: 5}},
	})
	if err != ErrCheckinNotFound {
		t.Fatalf("expected ErrCheckinNotFound for wrong user, got: %v", err)
	}
}
