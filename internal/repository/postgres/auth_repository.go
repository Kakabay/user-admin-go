package repository

import (
	"database/sql"
	"log/slog"
	"time"
	"user-admin/internal/domain"
	"user-admin/pkg/lib/utils"

	"user-admin/internal/config"

	"github.com/dgrijalva/jwt-go"
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

	err := row.Scan(&admin.ID, &admin.Username, &admin.Password, &admin.Roles)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrAdminNotFound
		}

		slog.Error("Error getting admin by username: %v", err)
		return nil, err
	}

	return &admin, nil
}

func (r *PostgresAdminAuthRepository) GenerateJWT(admin *domain.Admin) (string, error) {
	claims := jwt.MapClaims{
		"id": admin.ID,
		"role": admin.Roles,
		"exp": time.Now().Add(30 * time.Minute).Unix(), // Token expiration time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(r.JWTConfig.SecretKey))
	if err != nil {
		slog.Error("Error generating JWT: %v", utils.Err(err))
		return "", err
	}

	return tokenString, nil
}
