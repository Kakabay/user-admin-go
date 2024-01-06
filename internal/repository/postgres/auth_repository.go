package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
	"user-admin/internal/domain"
	"user-admin/pkg/lib/utils"

	"user-admin/internal/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type PostgresAdminAuthRepository struct {
	DB        *sql.DB
	JWTConfig config.JWT
}

func NewPostgresAdminAuthRepository(db *sql.DB, jwtConfig config.JWT) *PostgresAdminAuthRepository {
	return &PostgresAdminAuthRepository{DB: db, JWTConfig: jwtConfig}
}

func (r *PostgresAdminAuthRepository) GetAdminByUsername(username string) (*domain.Admin, error) {
	query := `
		SELECT id, username, password, role
		FROM admins
		WHERE username = $1
		LIMIT 1
	`

	row := r.DB.QueryRow(query, username)

	var admin domain.Admin

	err := row.Scan(&admin.ID, &admin.Username, &admin.Password, &admin.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("Admin was not found")
			return nil, domain.ErrAdminNotFound
		}

		slog.Error("Error getting admin by username: %v", err)
		return nil, err
	}

	return &admin, nil
}

func (r *PostgresAdminAuthRepository) GenerateAccessToken(admin *domain.Admin) (string, error) {
	claims := jwt.MapClaims{
		"id":   admin.ID,
		"role": admin.Role,
		"exp":  time.Now().Add(30 * time.Minute).Unix(), // Token expiration time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(r.JWTConfig.AccessSecretKey))
	if err != nil {
		slog.Error("Error generating access token: %v", utils.Err(err))
		return "", err
	}

	return tokenString, nil
}

func (r *PostgresAdminAuthRepository) GenerateRefreshToken(admin *domain.Admin) (string, error) {
	refreshTokenID := uuid.New().String()

	refreshClaims := jwt.MapClaims{
		"id":      refreshTokenID,
		"adminID": admin.ID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	refreshTokenString, err := refreshToken.SignedString([]byte(r.JWTConfig.RefreshSecretKey))
	if err != nil {
		return "", err
	}

	query := `
		UPDATE admins
		SET refresh_token = $1,
			refresh_token_created_at = TO_TIMESTAMP($2),
			refresh_token_expiration_time = TO_TIMESTAMP($3)
		WHERE id = $4
	`

	_, err = r.DB.Exec(query, refreshTokenString, refreshClaims["iat"].(int64), refreshClaims["exp"].(int64), admin.ID)
	if err != nil {
		slog.Error("Failed to update refresh token in database", utils.Err(err))
		return "", err
	}

	return refreshTokenString, nil
}

func (r *PostgresAdminAuthRepository) ValidateRefreshToken(refreshToken string) (map[string]interface{}, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(r.JWTConfig.RefreshSecretKey), nil
	})

	if err != nil || !token.Valid {
		slog.Error("Refresh token validation error: %v", err)
		return nil, fmt.Errorf("refresh token validation error: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims == nil {
		slog.Error("Invalid refresh token claims")
		return nil, fmt.Errorf("invalid refresh token claims")
	}

	return claims, nil
}

func (r *PostgresAdminAuthRepository) GetAdminByID(adminID int) (*domain.Admin, error) {
	query := `
		SELECT id, username, password, role
		FROM admins
		WHERE id = $1
		LIMIT 1
	`

	row := r.DB.QueryRow(query, adminID)

	var admin domain.Admin

	err := row.Scan(&admin.ID, &admin.Username, &admin.Password, &admin.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("Admin not found")
			return nil, domain.ErrAdminNotFound
		}

		slog.Error("Error getting admin by ID: %v", err)
		return nil, err
	}

	return &admin, nil
}

func (r *PostgresAdminAuthRepository) DeleteRefreshToken(refreshToken string) error {
	query := `
		DELETE FROM admins
		WHERE token = $1
	`

	_, err := r.DB.Exec(query, refreshToken)
	if err != nil {
		slog.Error("Failed to delete refresh token from database", utils.Err(err))
		return err
	}

	return nil
}
