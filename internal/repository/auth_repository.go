package repository

import (
	"user-admin/internal/domain"
)

type AdminAuthRepository interface {
	GetAdminByUsername(username string) (*domain.Admin, error)
	GenerateAccessToken(admin *domain.Admin) (string, error)
	GenerateRefreshToken(admin *domain.Admin) (string, error)
    ValidateRefreshToken(refreshToken string) (map[string]interface{}, error)
    GetAdminByID(adminID int) (*domain.Admin, error)
    DeleteRefreshToken(refreshToken string) error
}