package service

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
)

type Admin interface {
}

type adminService struct {
	adminRepository repository.Admin
	cache           repository.Cache
}

func NewAdminService(adminRepository repository.Admin, cache repository.Cache) Admin {
	return &adminService{
		adminRepository: adminRepository,
		cache:           cache,
	}
}

func (service *adminService) Register(ctx context.Context, email, pass string, roleType entity.AdminRoleType) (*entity.Admin, error) {
	return nil, nil
}
