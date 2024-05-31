package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/Daka-0424/my-go-server/pkg/controller/crypto"
	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const (
	ErrAlreadyHasEmail = "already has email"
	ErrAdminNotFound   = "admin not found"
	ErrEncryptPass     = "パスワード暗号化にエラーが発生しました。："
)

type IAdmin interface {
	Register(ctx context.Context, email, pass string, role entity.AdminRoleType) (*entity.Admin, error)
}

type adminUsecase struct {
	adminRepository repository.IAdmin
	transaction     repository.ITransaction
	localizer       *i18n.Localizer
}

func NewAdminUsecase(adminRepository repository.IAdmin, transaction repository.ITransaction, localizer *i18n.Localizer) IAdmin {
	return &adminUsecase{
		adminRepository: adminRepository,
		transaction:     transaction,
		localizer:       localizer,
	}
}

func (usecase *adminUsecase) Register(ctx context.Context, email, pass string, role entity.AdminRoleType) (*entity.Admin, error) {
	value, err := usecase.transaction.DoInTx(ctx, func(ctx context.Context) (interface{}, error) {
		exsist := usecase.adminRepository.Exsists(ctx, email)
		if exsist {
			return nil, errors.New(ErrAlreadyHasEmail)
		}

		encryptPw, err := crypto.PasswordEncrypt(pass)
		if err != nil {
			fmt.Println(ErrEncryptPass, err)
			return nil, err
		}

		count, err := usecase.adminRepository.CountAll(ctx)
		if err != nil {
			return nil, err
		}

		if count == 0 {
			role = entity.AdminRoleTypeMaster
		}

		admin, err := usecase.adminRepository.Register(ctx, email, encryptPw, role)
		if err != nil {
			return nil, err
		}

		return admin, nil
	})

	if err != nil {
		return nil, err
	}

	return value.(*entity.Admin), nil
}
