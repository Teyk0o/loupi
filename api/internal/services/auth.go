// Package services contains the business logic layer.
package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/teyk0o/loupi/api/internal/config"
	"github.com/teyk0o/loupi/api/internal/models"
)

// Common errors returned by the auth service.
var (
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidToken       = errors.New("invalid or expired token")
)

const (
	accessTokenDuration  = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour
)

// AuthService handles authentication business logic.
type AuthService struct {
	db  *pgxpool.Pool
	cfg *config.Config
}

// NewAuthService creates a new AuthService instance.
func NewAuthService(db *pgxpool.Pool, cfg *config.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

// Register creates a new user account with email and password.
func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) (*models.TokenResponse, error) {
	var exists bool
	err := s.db.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}
	if exists {
		return nil, ErrEmailAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), s.cfg.BcryptCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	var user models.User
	hashStr := string(hash)
	err = s.db.QueryRow(ctx,
		`INSERT INTO users (email, password_hash, first_name) VALUES ($1, $2, $3)
		 RETURNING id, email, first_name, password_hash, oauth_provider, oauth_id, email_verified, created_at, updated_at`,
		req.Email, hashStr, req.FirstName,
	).Scan(&user.ID, &user.Email, &user.FirstName, &user.PasswordHash, &user.OAuthProvider, &user.OAuthID, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return s.generateTokens(&user)
}

// Login authenticates a user with email and password.
func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.TokenResponse, error) {
	user, err := s.getUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if user.PasswordHash == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.generateTokens(user)
}

// RefreshToken generates a new access token from a valid refresh token.
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*models.TokenResponse, error) {
	claims, err := s.validateToken(refreshToken, "refresh")
	if err != nil {
		return nil, ErrInvalidToken
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, ErrInvalidToken
	}

	user, err := s.GetUserByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	return s.generateTokens(user)
}

// GetUserByID retrieves a user by their ID.
func (s *AuthService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(ctx,
		`SELECT id, email, first_name, password_hash, oauth_provider, oauth_id, email_verified, created_at, updated_at
		 FROM users WHERE id = $1`,
		id,
	).Scan(&user.ID, &user.Email, &user.FirstName, &user.PasswordHash, &user.OAuthProvider, &user.OAuthID, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// DeleteAccount removes a user and all associated data.
func (s *AuthService) DeleteAccount(ctx context.Context, userID uuid.UUID) error {
	result, err := s.db.Exec(ctx, "DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}
	return nil
}

// getUserByEmail retrieves a user by their email address.
func (s *AuthService) getUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(ctx,
		`SELECT id, email, first_name, password_hash, oauth_provider, oauth_id, email_verified, created_at, updated_at
		 FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Email, &user.FirstName, &user.PasswordHash, &user.OAuthProvider, &user.OAuthID, &user.EmailVerified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// generateTokens creates a new pair of access and refresh tokens for a user.
func (s *AuthService) generateTokens(user *models.User) (*models.TokenResponse, error) {
	now := time.Now()

	accessToken, err := s.createToken(user.ID, "access", now.Add(accessTokenDuration))
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	refreshToken, err := s.createToken(user.ID, "refresh", now.Add(refreshTokenDuration))
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	return &models.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(accessTokenDuration.Seconds()),
		User:         user.ToResponse(),
	}, nil
}

// createToken generates a signed JWT with the given claims.
func (s *AuthService) createToken(userID uuid.UUID, tokenType string, expiresAt time.Time) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		Audience:  jwt.ClaimStrings{tokenType},
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		Issuer:    "loupi",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.JWTSecret))
}

// validateToken verifies a JWT and returns its claims.
func (s *AuthService) validateToken(tokenString string, expectedType string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	audiences, err := claims.GetAudience()
	if err != nil {
		return nil, ErrInvalidToken
	}

	found := false
	for _, aud := range audiences {
		if aud == expectedType {
			found = true
			break
		}
	}
	if !found {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ValidateAccessToken validates an access token and returns the user ID.
func (s *AuthService) ValidateAccessToken(tokenString string) (uuid.UUID, error) {
	claims, err := s.validateToken(tokenString, "access")
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, ErrInvalidToken
	}

	return userID, nil
}
