package service

import (
	"context"

	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
)

type IAdmin interface {
}

type adminService struct {
	adminRepository repository.IAdmin
	cache           repository.ICache
}

func NewAdminService(adminRepository repository.IAdmin, cache repository.ICache) IAdmin {
	return &adminService{
		adminRepository: adminRepository,
		cache:           cache,
	}
}

func (service *adminService) Register(ctx context.Context, email, pass string, roleType entity.AdminRoleType) (*entity.Admin, error) {
	return nil, nil
}
