package usecase

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/controller/crypto"
	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/gomail.v2"
)

const (
	EMailSubject       = "EMAIL_SUBJECT"
	EMailBody          = "EMAIL_BODY"
	ErrAlreadyHasEmail = "already has email"
	ErrCreateRedisKey  = "ランダムな文字列の生成に失敗しました。"
	ErrAdminNotFound   = "admin not found"
	ErrEncryptPass     = "パスワード暗号化にエラーが発生しました。："
)

type IAdmin interface {
	TempRegister(ctx context.Context, email string) error
	Register(ctx context.Context, email, pass string, role entity.AdminRoleType) (*entity.Admin, error)
}

type adminUsecase struct {
	adminRepository repository.IAdmin
	transaction     repository.ITransaction
	cache           repository.ICache
	cfg             *config.Config
	localizer       *i18n.Localizer
}

func NewAdminUsecase(
	adminRepository repository.IAdmin,
	transaction repository.ITransaction,
	cache repository.ICache,
	cfg *config.Config,
	localizer *i18n.Localizer,
) IAdmin {
	return &adminUsecase{
		adminRepository: adminRepository,
		transaction:     transaction,
		cache:           cache,
		cfg:             cfg,
		localizer:       localizer,
	}
}

func (usecase *adminUsecase) TempRegister(ctx context.Context, email string) error {
	// 登録済みE-Mail
	exsists := usecase.adminRepository.Exsists(ctx, email)
	if exsists {
		return errors.New(ErrAlreadyHasEmail)
	}

	// ランダムな文字列生成
	key := make([]byte, 64)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return errors.New(ErrCreateRedisKey)
	}
	redisKey := base64.URLEncoding.EncodeToString(key)

	// Redisに保存
	usecase.cache.Set(ctx, redisKey, []byte(email), time.Hour)

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", usecase.cfg.Admin.RegisterEmailSender)

	// Set E-Mail receiver
	m.SetHeader("To", email)

	m.SetHeader("Bcc", usecase.cfg.Admin.RegisterEmailSender)

	// Set E-Mail subject
	c := &i18n.LocalizeConfig{MessageID: EMailSubject}
	m.SetHeader("Subject", usecase.localizer.MustLocalize(c))

	// Set E-Mail body. You can set plain text or html with text/html
	c = &i18n.LocalizeConfig{MessageID: EMailBody}
	body := usecase.localizer.MustLocalize(c)
	domain := usecase.cfg.Settings.BaseDomain + "/admin/register/" + redisKey
	m.SetBody("text/plain", fmt.Sprintf(body, domain))

	// Setting for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, usecase.cfg.Admin.RegisterEmailSender, usecase.cfg.Admin.RegisterEmailPass)

	// TLS configuration
	d.TLSConfig = &tls.Config{InsecureSkipVerify: usecase.cfg.IsDevelopment()}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
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
