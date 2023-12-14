package service

import (
	"log/slog"
	"user-admin/internal/domain"
	"user-admin/internal/repository"
	"user-admin/pkg/lib/utils"

	"golang.org/x/crypto/bcrypt"
)

type AdminAuthService struct {
	AdminAuthRepository repository.AdminAuthRepository
}

func NewAdminAuthService(adminAuthRepository repository.AdminAuthRepository) *AdminAuthService {
	return &AdminAuthService{AdminAuthRepository: adminAuthRepository}
}

func (s *AdminAuthService) LoginAdmin(username, password string) (string, error) {
    admin, err := s.AdminAuthRepository.GetAdminByUsername(username)
    if err != nil {
        slog.Error("Error getting admin by username:", utils.Err(err))
        return "", err
    }

    err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
    if err != nil {
        slog.Error("Error comparing passwords:", utils.Err(err))
        return "", domain.ErrInvalidCredentials
    }

    token, err := s.AdminAuthRepository.GenerateJWT(admin)
    if err != nil {
        slog.Error("Error generating JWT:", utils.Err(err))
        return "", err
    }

    return token, nil
}
