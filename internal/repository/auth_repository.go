package repository

import "user-admin/internal/domain"

type AdminAuthRepository interface {
	// Methods for authentication
	GetAdminByUsername(username string) (*domain.Admin, error)
	GenerateAccessToken(admin *domain.Admin) (string, error)
}