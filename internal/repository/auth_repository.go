package repository

import "user-admin/internal/domain"

type AdminAuthRepository interface {
	// Methods for authentication
	GetAdminByUsername(username string) (*domain.Admin, error)
	GenerateJWT(adminID int32) (string, error)
	ValidateJWT(token string) (int32, error)
}