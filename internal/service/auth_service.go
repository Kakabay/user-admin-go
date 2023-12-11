package service

import (
	"user-admin/internal/domain"
	"user-admin/internal/repository"

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
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	token, err := s.AdminAuthRepository.GenerateJWT(admin.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}