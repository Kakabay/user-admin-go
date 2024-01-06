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

func (s *AdminAuthService) LoginAdmin(username, password string) (string, string, error) {
	admin, err := s.AdminAuthRepository.GetAdminByUsername(username)
	if err != nil {
		slog.Error("Error getting admin by username:", utils.Err(err))
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		slog.Error("Error comparing passwords:", utils.Err(err))
		return "", "", domain.ErrInvalidCredentials
	}

	accessToken, err := s.AdminAuthRepository.GenerateAccessToken(admin)
	if err != nil {
		slog.Error("Error generating access token:", utils.Err(err))
		return "", "", err
	}

	refreshToken, err := s.AdminAuthRepository.GenerateRefreshToken(admin)
	if err != nil {
		slog.Error("Error generating refresh token:", utils.Err(err))
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AdminAuthService) RefreshTokens(refreshToken string) (string, error) {
	claims, err := s.AdminAuthRepository.ValidateRefreshToken(refreshToken)
	if err != nil {
		slog.Error("Error validating refresh token:", utils.Err(err))
		return "", err
	}

	adminIDFloat, ok := claims["adminID"].(float64)
	if !ok {
		slog.Error("AdminID not found or not a number in refresh token claims")
		return "", domain.ErrInvalidRefreshToken
	}

	// Convert adminID to int
	adminID := int(adminIDFloat)

	admin, err := s.AdminAuthRepository.GetAdminByID(adminID)
	if err != nil {
		slog.Error("Error getting admin by ID:", utils.Err(err))
		return "", err
	}

	newAccessToken, err := s.AdminAuthRepository.GenerateAccessToken(admin)
	if err != nil {
		slog.Error("Error generating access token:", utils.Err(err))
		return "", err
	}

	return newAccessToken, nil
}
