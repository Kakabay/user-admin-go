package service

import (
	"user-admin/internal/domain"
	"user-admin/internal/repository"
)

type UserService struct {
	UserRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{UserRepository: userRepository}
}

func (s *UserService) GetAllUsers(page, pageSize int) (*domain.UsersList, error) {
	return s.UserRepository.GetAllUsers(page, pageSize)
}

func (s *UserService) GetUserByID(id int32) (*domain.GetUserResponse, error) {
	return s.UserRepository.GetUserByID(id)
}

func (s *UserService) CreateUser(request *domain.CreateUserRequest) (*domain.CreateUserResponse, error) {
	return s.UserRepository.CreateUser(request)
}

func (s *UserService) UpdateUser(request *domain.UpdateUserRequest) (*domain.UpdateUserResponse, error) {
	return s.UserRepository.UpdateUser(request)
}

func (s *UserService) DeleteUser(id int32) error {
	return s.UserRepository.DeleteUser(id)
}

func (s *UserService) BlockUser(id int32) error {
	return s.UserRepository.BlockUser(id)
}

func (s *UserService) UnblockUser(id int32) error {
	return s.UserRepository.UnblockUser(id)
}

func (s *UserService) SearchUsers(query string, page, pageSize int) (*domain.UsersList, error) {
	return s.UserRepository.SearchUsers(query, page, pageSize)
}