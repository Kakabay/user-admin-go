package service

import (
	"user-admin/internal/domain"
	"user-admin/internal/repository"
)

type AdminService struct {
	AdminRepository repository.AdminRepository
}

func NewAdminService(adminRepository repository.AdminRepository) *AdminService {
	return &AdminService{AdminRepository: adminRepository}
}

/*
	func (s *AdminService) GetAllAdmins(page, pageSize int) (*domain.AdminsList, error) {
		return s.AdminRepository.GetAllAdmins(page, pageSize)
	}
*/
func (s *AdminService) GetAdminByID(id int32) (*domain.CommonAdminResponse, error) {
	return s.AdminRepository.GetAdminByID(id)
}

func (s *AdminService) CreateAdmin(request *domain.CreateAdminRequest) (*domain.CommonAdminResponse, error) {
	return s.AdminRepository.CreateAdmin(request)
}

/*
func (s *AdminService) UpdateAdmin(request *domain.Admin) (*domain.Admin, error) {
	return s.AdminRepository.UpdateAdmin(request)
}

func (s *AdminService) DeleteAdmin(id int32) error {
	return s.AdminRepository.DeleteAdmin(id)
}

func (s *AdminService) SearchAdmins(query string, page, pageSize int) (*domain.AdminsList, error) {
	return s.AdminRepository.SearchAdmins(query, page, pageSize)
}
*/
