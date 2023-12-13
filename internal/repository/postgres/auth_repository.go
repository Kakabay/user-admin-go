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
		SELECT id, username, password
		FROM admins
		WHERE username = $1
		LIMIT 1
	`

	row := r.DB.QueryRow(query, username)

	var admin domain.Admin

	err := row.Scan(&admin.ID, &admin.Username, &admin.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrAdminNotFound
		}

		slog.Error("Error getting admin by username: %v", err)
		return nil, err
	}

	return &admin, nil
}

func (r *PostgresAdminAuthRepository) GenerateJWT(adminID int32) (string, error) {
	claims := jwt.MapClaims{
		"id": adminID,
		"exp": time.Now().Add(time.Minute * 30).Unix(), // Token expiration time
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(r.JWTConfig.SecretKey))
	if err != nil {
		slog.Error("Error generating JWT: %v", utils.Err(err))
		return "", err
	}

	return tokenString, nil
}

func (r *PostgresAdminAuthRepository) ValidateJWT(tokenString string) (int32, error) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(r.JWTConfig.SecretKey), nil
    })

    if err != nil {
		slog.Error("Error validating JWT: %v", err)
        return 0, err
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        adminID, ok := claims["id"].(float64)
        if !ok {
			slog.Error("invalid ID claim in JWT")
            return 0, fmt.Errorf("invalid ID claim in JWT")
        }

        return int32(adminID), nil
    }

	slog.Error("invalid JWT token")
    return 0, fmt.Errorf("invalid JWT token")
}
