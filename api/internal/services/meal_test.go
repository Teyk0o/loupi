package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/teyk0o/loupi/api/internal/models"
)

// createTestUser registers a test user and returns their ID.
func createTestUser(t *testing.T, pool interface{ QueryRow(ctx context.Context, sql string, args ...interface{}) interface{ Scan(dest ...interface{}) error } }, email string) uuid.UUID {
	t.Helper()
	// Use the auth service for convenience
	return uuid.Nil // placeholder, overridden below
}

// setupTestMealEnv creates a test user and returns the user ID and meal service.
func setupTestMealEnv(t *testing.T) (*MealService, uuid.UUID, func()) {
	t.Helper()
	pool := setupTestDB(t)
	rdb := setupTestRedis(t)
	cfg := testConfig()
	enc := testEncryptor(t)
	authSvc := NewAuthService(pool, rdb, cfg)
	mealSvc := NewMealService(pool, enc)
	ctx := context.Background()

	email := "test-meal-" + uuid.New().String()[:8] + "@loupi.test"
	resp, err := authSvc.Register(ctx, models.RegisterRequest{Email: email, Password: "SecurePass123#"})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	cleanup := func() {
		_ = authSvc.DeleteAccount(ctx, resp.User.ID)
		rdb.Close()
		pool.Close()
	}

	return mealSvc, resp.User.ID, cleanup
}

func TestMeal_Create(t *testing.T) {
	svc, userID, cleanup := setupTestMealEnv(t)
	defer cleanup()

	meal, err := svc.Create(context.Background(), userID, models.CreateMealRequest{
		Description: "Pasta with tomato sauce",
		Category:    "homemade",
		MealTime:    time.Now().Format(time.RFC3339),
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if meal.Description != "Pasta with tomato sauce" {
		t.Errorf("expected description 'Pasta with tomato sauce', got '%s'", meal.Description)
	}
	if meal.Category != "homemade" {
		t.Errorf("expected category 'homemade', got '%s'", meal.Category)
	}
	if meal.UserID != userID {
		t.Errorf("expected user ID %s, got %s", userID, meal.UserID)
	}
}

func TestMeal_Create_InvalidTime(t *testing.T) {
	svc, userID, cleanup := setupTestMealEnv(t)
	defer cleanup()

	_, err := svc.Create(context.Background(), userID, models.CreateMealRequest{
		Description: "Test",
		Category:    "snack",
		MealTime:    "not-a-date",
	})
	if err == nil {
		t.Fatal("expected error for invalid meal_time")
	}
}

func TestMeal_GetByID(t *testing.T) {
	svc, userID, cleanup := setupTestMealEnv(t)
	defer cleanup()

	created, err := svc.Create(context.Background(), userID, models.CreateMealRequest{
		Description: "Salad",
		Category:    "restaurant",
		MealTime:    time.Now().Format(time.RFC3339),
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	meal, err := svc.GetByID(context.Background(), userID, created.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if meal.ID != created.ID {
		t.Errorf("expected ID %s, got %s", created.ID, meal.ID)
	}
}

func TestMeal_GetByID_NotFound(t *testing.T) {
	svc, userID, cleanup := setupTestMealEnv(t)
	defer cleanup()

	_, err := svc.GetByID(context.Background(), userID, uuid.New())
	if err != ErrMealNotFound {
		t.Fatalf("expected ErrMealNotFound, got: %v", err)
	}
}

func TestMeal_GetByID_WrongUser(t *testing.T) {
	svc, userID, cleanup := setupTestMealEnv(t)
	defer cleanup()

	created, err := svc.Create(context.Background(), userID, models.CreateMealRequest{
		Description: "Private meal",
		Category:    "homemade",
		MealTime:    time.Now().Format(time.RFC3339),
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	// Another user should not see this meal
	_, err = svc.GetByID(context.Background(), uuid.New(), created.ID)
	if err != ErrMealNotFound {
		t.Fatalf("expected ErrMealNotFound for wrong user, got: %v", err)
	}
}

func TestMeal_ListByDate(t *testing.T) {
	svc, userID, cleanup := setupTestMealEnv(t)
	defer cleanup()

	today := time.Now().Format("2006-01-02")
	now := time.Now()

	_, _ = svc.Create(context.Background(), userID, models.CreateMealRequest{
		Description: "Breakfast",
		Category:    "homemade",
		MealTime:    now.Format(time.RFC3339),
	})
	_, _ = svc.Create(context.Background(), userID, models.CreateMealRequest{
		Description: "Lunch",
		Category:    "restaurant",
		MealTime:    now.Format(time.RFC3339),
	})

	meals, err := svc.ListByDate(context.Background(), userID, today)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(meals) < 2 {
		t.Errorf("expected at least 2 meals, got %d", len(meals))
	}
}

func TestMeal_ListByDate_Empty(t *testing.T) {
	svc, userID, cleanup := setupTestMealEnv(t)
	defer cleanup()

	meals, err := svc.ListByDate(context.Background(), userID, "2020-01-01")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(meals) != 0 {
		t.Errorf("expected 0 meals, got %d", len(meals))
	}
}

func TestMeal_Update(t *testing.T) {
	svc, userID, cleanup := setupTestMealEnv(t)
	defer cleanup()

	created, err := svc.Create(context.Background(), userID, models.CreateMealRequest{
		Description: "Old description",
		Category:    "homemade",
		MealTime:    time.Now().Format(time.RFC3339),
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	newDesc := "Updated description"
	updated, err := svc.Update(context.Background(), userID, created.ID, models.UpdateMealRequest{
		Description: &newDesc,
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if updated.Description != newDesc {
		t.Errorf("expected description '%s', got '%s'", newDesc, updated.Description)
	}
	if updated.Category != "homemade" {
		t.Errorf("expected category to remain 'homemade', got '%s'", updated.Category)
	}
}

func TestMeal_Delete(t *testing.T) {
	svc, userID, cleanup := setupTestMealEnv(t)
	defer cleanup()

	created, err := svc.Create(context.Background(), userID, models.CreateMealRequest{
		Description: "To be deleted",
		Category:    "snack",
		MealTime:    time.Now().Format(time.RFC3339),
	})
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	err = svc.Delete(context.Background(), userID, created.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	_, err = svc.GetByID(context.Background(), userID, created.ID)
	if err != ErrMealNotFound {
		t.Fatalf("expected ErrMealNotFound after delete, got: %v", err)
	}
}

func TestMeal_Delete_NotFound(t *testing.T) {
	svc, userID, cleanup := setupTestMealEnv(t)
	defer cleanup()

	err := svc.Delete(context.Background(), userID, uuid.New())
	if err != ErrMealNotFound {
		t.Fatalf("expected ErrMealNotFound, got: %v", err)
	}
}

func TestMeal_Checkin_CreateAndList(t *testing.T) {
	svc, userID, cleanup := setupTestMealEnv(t)
	defer cleanup()

	meal, err := svc.Create(context.Background(), userID, models.CreateMealRequest{
		Description: "Meal with checkin",
		Category:    "homemade",
		MealTime:    time.Now().Format(time.RFC3339),
	})
	if err != nil {
		t.Fatalf("create meal failed: %v", err)
	}

	checkin, err := svc.CreateCheckin(context.Background(), userID, meal.ID, models.CreateCheckinRequest{
		DelayHours: 6,
		Symptoms:   []models.SymptomDetail{{Type: "bloating", Severity: 3}},
	})
	if err != nil {
		t.Fatalf("create checkin failed: %v", err)
	}
	if checkin.MealID != meal.ID {
		t.Errorf("expected meal_id %s, got %s", meal.ID, checkin.MealID)
	}
	if checkin.DelayHours != 6 {
		t.Errorf("expected delay_hours 6, got %d", checkin.DelayHours)
	}

	checkins, err := svc.GetCheckins(context.Background(), userID, meal.ID)
	if err != nil {
		t.Fatalf("list checkins failed: %v", err)
	}
	if len(checkins) != 1 {
		t.Errorf("expected 1 checkin, got %d", len(checkins))
	}
}

func TestMeal_Checkin_WrongMeal(t *testing.T) {
	svc, userID, cleanup := setupTestMealEnv(t)
	defer cleanup()

	_, err := svc.CreateCheckin(context.Background(), userID, uuid.New(), models.CreateCheckinRequest{
		DelayHours: 8,
		Symptoms:   []models.SymptomDetail{{Type: "nausea", Severity: 2}},
	})
	if err != ErrMealNotFound {
		t.Fatalf("expected ErrMealNotFound, got: %v", err)
	}
}
