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

func (s *UserService) GetAllUsers() (*domain.UsersList, error) {
	return s.UserRepository.GetAllUsers()
}

func (s *UserService) GetUserByID(id int32) (*domain.CommonUserResponse, error) {
	return s.UserRepository.GetUserByID()
}

func (s *UserService) CreateUser(request *domain.CreateUserRequest) (*domain.CommonUserResponse, error) {
	return s.UserRepository.CreateUser()
}

func (s *UserService) UpdateUser(request *domain.UpdateUserRequest) (*domain.CommonUserResponse, error) {
	return s.UserRepository.UpdateUser()
}

func (s *UserService) DeleteUser(id int32) error {
	return s.UserRepository.DeleteUser()
}

func (s *UserService) BlockUser(id int32) error {
	return s.UserRepository.BlockUser()
}

func (s *UserService) UnblockUser(id int32) error {
	return s.UserRepository.UnblockUser()
}