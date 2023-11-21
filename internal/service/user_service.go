package service

import (
	"user-admin/internal/domain"
	"user-admin/internal/repository"
)

type UserService struct {
	UserRepository repository.UserRepository
}

func (s *UserService) GetAllUsers() (*domain.UsersList, error) {
	return s.UserRepository.GetAllUsers()
}