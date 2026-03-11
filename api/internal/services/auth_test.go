package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/teyk0o/loupi/api/internal/config"
	"github.com/teyk0o/loupi/api/internal/models"
	"github.com/teyk0o/loupi/api/internal/utils"
)

// testConfig returns a config suitable for testing.
func testConfig() *config.Config {
	return &config.Config{
		Env:           "development",
		JWTSecret:     "test-secret-key-for-unit-tests-only",
		EncryptionKey: "0000000000000000000000000000000000000000000000000000000000000000",
		BcryptCost:    4, // Low cost for fast tests
	}
}

// setupTestDB connects to the test database.
// Requires PostgreSQL running (docker-compose.dev.yml).
func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dbURL := "postgres://loupi:loupi_dev@localhost:5432/loupi?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		t.Skipf("Skipping test: database not available: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		t.Skipf("Skipping test: database not reachable: %v", err)
	}

	return pool
}

// setupTestRedis connects to the test Redis instance.
// Requires Redis running (docker-compose.dev.yml).
func setupTestRedis(t *testing.T) *redis.Client {
	t.Helper()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "redis_dev",
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		rdb.Close()
		t.Skipf("Skipping test: Redis not available: %v", err)
	}

	return rdb
}

// testEncryptor returns an Encryptor suitable for testing.
func testEncryptor(t *testing.T) *utils.Encryptor {
	t.Helper()
	enc, err := utils.NewEncryptor("0000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		t.Fatalf("failed to create test encryptor: %v", err)
	}
	return enc
}

// cleanupTestUser removes a test user by email.
func cleanupTestUser(t *testing.T, pool *pgxpool.Pool, email string) {
	t.Helper()
	_, _ = pool.Exec(context.Background(), "DELETE FROM users WHERE email = $1", email)
}

func TestRegister_Success(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	rdb := setupTestRedis(t)
	defer rdb.Close()

	cfg := testConfig()
	svc := NewAuthService(pool, rdb, cfg)
	ctx := context.Background()
	email := "test-register-success@loupi.test"

	cleanupTestUser(t, pool, email)
	defer cleanupTestUser(t, pool, email)

	resp, err := svc.Register(ctx, models.RegisterRequest{
		Email:    email,
		Password: "SecurePass123#",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if resp.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
	if resp.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}
	if resp.User.Email != email {
		t.Errorf("expected email %s, got %s", email, resp.User.Email)
	}
	if resp.User.EmailVerified {
		t.Error("expected email_verified to be false")
	}
	if resp.ExpiresIn != 900 {
		t.Errorf("expected expires_in 900, got %d", resp.ExpiresIn)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	rdb := setupTestRedis(t)
	defer rdb.Close()

	cfg := testConfig()
	svc := NewAuthService(pool, rdb, cfg)
	ctx := context.Background()
	email := "test-register-dup@loupi.test"

	cleanupTestUser(t, pool, email)
	defer cleanupTestUser(t, pool, email)

	_, err := svc.Register(ctx, models.RegisterRequest{
		Email:    email,
		Password: "SecurePass123#",
	})
	if err != nil {
		t.Fatalf("first register failed: %v", err)
	}

	_, err = svc.Register(ctx, models.RegisterRequest{
		Email:    email,
		Password: "AnotherPass456#",
	})
	if err != ErrEmailAlreadyExists {
		t.Fatalf("expected ErrEmailAlreadyExists, got: %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	rdb := setupTestRedis(t)
	defer rdb.Close()

	cfg := testConfig()
	svc := NewAuthService(pool, rdb, cfg)
	ctx := context.Background()
	email := "test-login-success@loupi.test"
	password := "SecurePass123#"

	cleanupTestUser(t, pool, email)
	defer cleanupTestUser(t, pool, email)

	_, err := svc.Register(ctx, models.RegisterRequest{Email: email, Password: password})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	resp, err := svc.Login(ctx, models.LoginRequest{Email: email, Password: password})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if resp.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
	if resp.User.Email != email {
		t.Errorf("expected email %s, got %s", email, resp.User.Email)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	rdb := setupTestRedis(t)
	defer rdb.Close()

	cfg := testConfig()
	svc := NewAuthService(pool, rdb, cfg)
	ctx := context.Background()
	email := "test-login-wrong@loupi.test"

	cleanupTestUser(t, pool, email)
	defer cleanupTestUser(t, pool, email)

	_, err := svc.Register(ctx, models.RegisterRequest{Email: email, Password: "SecurePass123#"})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	_, err = svc.Login(ctx, models.LoginRequest{Email: email, Password: "WrongPassword"})
	if err != ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got: %v", err)
	}
}

func TestLogin_NonexistentUser(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	rdb := setupTestRedis(t)
	defer rdb.Close()

	cfg := testConfig()
	svc := NewAuthService(pool, rdb, cfg)
	ctx := context.Background()

	_, err := svc.Login(ctx, models.LoginRequest{Email: "nonexistent@loupi.test", Password: "whatever"})
	if err != ErrInvalidCredentials {
		t.Fatalf("expected ErrInvalidCredentials, got: %v", err)
	}
}

func TestRefreshToken_Success(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	rdb := setupTestRedis(t)
	defer rdb.Close()

	cfg := testConfig()
	svc := NewAuthService(pool, rdb, cfg)
	ctx := context.Background()
	email := "test-refresh@loupi.test"

	cleanupTestUser(t, pool, email)
	defer cleanupTestUser(t, pool, email)

	resp, err := svc.Register(ctx, models.RegisterRequest{Email: email, Password: "SecurePass123#"})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	newResp, err := svc.RefreshToken(ctx, resp.RefreshToken)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if newResp.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
	if newResp.User.Email != email {
		t.Errorf("expected email %s, got %s", email, newResp.User.Email)
	}
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	rdb := setupTestRedis(t)
	defer rdb.Close()

	cfg := testConfig()
	svc := NewAuthService(pool, rdb, cfg)
	ctx := context.Background()

	_, err := svc.RefreshToken(ctx, "invalid-token")
	if err != ErrInvalidToken {
		t.Fatalf("expected ErrInvalidToken, got: %v", err)
	}
}

func TestRefreshToken_AccessTokenRejected(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	rdb := setupTestRedis(t)
	defer rdb.Close()

	cfg := testConfig()
	svc := NewAuthService(pool, rdb, cfg)
	ctx := context.Background()
	email := "test-refresh-reject@loupi.test"

	cleanupTestUser(t, pool, email)
	defer cleanupTestUser(t, pool, email)

	resp, err := svc.Register(ctx, models.RegisterRequest{Email: email, Password: "SecurePass123#"})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// Using an access token as refresh token should fail
	_, err = svc.RefreshToken(ctx, resp.AccessToken)
	if err != ErrInvalidToken {
		t.Fatalf("expected ErrInvalidToken when using access token as refresh, got: %v", err)
	}
}

func TestValidateAccessToken_Success(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	rdb := setupTestRedis(t)
	defer rdb.Close()

	cfg := testConfig()
	svc := NewAuthService(pool, rdb, cfg)
	ctx := context.Background()
	email := "test-validate@loupi.test"

	cleanupTestUser(t, pool, email)
	defer cleanupTestUser(t, pool, email)

	resp, err := svc.Register(ctx, models.RegisterRequest{Email: email, Password: "SecurePass123#"})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	userID, _, err := svc.ValidateAccessToken(ctx, resp.AccessToken)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if userID != resp.User.ID {
		t.Errorf("expected user ID %s, got %s", resp.User.ID, userID)
	}
}

func TestValidateAccessToken_RefreshTokenRejected(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	rdb := setupTestRedis(t)
	defer rdb.Close()

	cfg := testConfig()
	svc := NewAuthService(pool, rdb, cfg)
	ctx := context.Background()
	email := "test-validate-reject@loupi.test"

	cleanupTestUser(t, pool, email)
	defer cleanupTestUser(t, pool, email)

	resp, err := svc.Register(ctx, models.RegisterRequest{Email: email, Password: "SecurePass123#"})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// Using a refresh token as access token should fail
	_, _, err = svc.ValidateAccessToken(ctx, resp.RefreshToken)
	if err == nil {
		t.Fatal("expected error when using refresh token as access token")
	}
}

func TestDeleteAccount_Success(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	rdb := setupTestRedis(t)
	defer rdb.Close()

	cfg := testConfig()
	svc := NewAuthService(pool, rdb, cfg)
	ctx := context.Background()
	email := "test-delete@loupi.test"

	cleanupTestUser(t, pool, email)

	resp, err := svc.Register(ctx, models.RegisterRequest{Email: email, Password: "SecurePass123#"})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	err = svc.DeleteAccount(ctx, resp.User.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	// Verify user is gone
	_, err = svc.GetUserByID(ctx, resp.User.ID)
	if err != ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
}

func TestDeleteAccount_NonexistentUser(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	rdb := setupTestRedis(t)
	defer rdb.Close()

	cfg := testConfig()
	svc := NewAuthService(pool, rdb, cfg)
	ctx := context.Background()

	err := svc.DeleteAccount(ctx, uuid.New())
	if err != ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
}

func TestGetUserByID_Success(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	rdb := setupTestRedis(t)
	defer rdb.Close()

	cfg := testConfig()
	svc := NewAuthService(pool, rdb, cfg)
	ctx := context.Background()
	email := "test-getbyid@loupi.test"

	cleanupTestUser(t, pool, email)
	defer cleanupTestUser(t, pool, email)

	resp, err := svc.Register(ctx, models.RegisterRequest{Email: email, Password: "SecurePass123#"})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	user, err := svc.GetUserByID(ctx, resp.User.ID)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if user.Email != email {
		t.Errorf("expected email %s, got %s", email, user.Email)
	}
	if user.PasswordHash == nil {
		t.Error("expected password hash to be set")
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	defer pool.Close()
	rdb := setupTestRedis(t)
	defer rdb.Close()

	cfg := testConfig()
	svc := NewAuthService(pool, rdb, cfg)
	ctx := context.Background()

	_, err := svc.GetUserByID(ctx, uuid.New())
	if err != ErrUserNotFound {
		t.Fatalf("expected ErrUserNotFound, got: %v", err)
	}
}
