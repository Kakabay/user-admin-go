package repository

import "user-admin/internal/domain"

type UserRepository interface {
	GetAllUsers(page, pageSize int) (*domain.UsersList, error)
	GetUserByID(id int32) (*domain.GetUserResponse, error)
	CreateUser(request *domain.CreateUserRequest) (*domain.CreateUserResponse, error)
	UpdateUser(request *domain.UpdateUserRequest) (*domain.UpdateUserResponse, error)
	DeleteUser(id int32) error
	BlockUser(id int32) error
	UnblockUser(id int32) error
}