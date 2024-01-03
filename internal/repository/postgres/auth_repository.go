package repository

import (
	"database/sql"
	"log/slog"
	"time"
	"user-admin/internal/domain"
	"user-admin/pkg/lib/utils"

	"user-admin/internal/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type PostgresAdminAuthRepository struct {
	DB *sql.DB
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
		"id": admin.ID,
		"role": admin.Role,
		"exp": time.Now().Add(30 * time.Minute).Unix(), // Token expiration time
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
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	refreshTokenString, err := refreshToken.SignedString([]byte(r.JWTConfig.RefreshSecretKey))
	if err != nil {
		return "", err
	}

	query := `
		INSERT INTO refresh_tokens (admin_id, token, expiration_time)
		VALUES ($1, $2, TO_TIMESTAMP($3))
	`
	_, err = r.DB.Exec(query, admin.ID, refreshTokenString, refreshClaims["exp"])
	if err != nil {
		slog.Error("Failed to save generated refresh tokens in database", utils.Err(err))
		return "", err
	}

	return refreshTokenString, nil
}
