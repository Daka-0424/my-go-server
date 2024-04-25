package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strconv"
	"time"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/controller/formatter"
	"github.com/Daka-0424/my-go-server/pkg/controller/middleware"
	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/Songmu/flextime"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Session interface {
	CreateSession(ctx context.Context, userId uint, uuid, device, appVersion string, platformNumber uint) (*model.Session, error)
}

type sessionUsecase struct {
	cfg            *config.Config
	localizer      *i18n.Localizer
	cache          repository.Cache
	transaction    repository.Transaction
	userRepository repository.User
}

func NewSessionUsecase(
	cfg *config.Config,
	lc *i18n.Localizer,
	cache repository.Cache,
	transaction repository.Transaction,
	userRepository repository.User) Session {
	return &sessionUsecase{
		cfg:            cfg,
		localizer:      lc,
		cache:          cache,
		transaction:    transaction,
		userRepository: userRepository,
	}
}

func (usecase *sessionUsecase) CreateSession(ctx context.Context, userId uint, uuid, device, appVersion string, platformNumber uint) (*model.Session, error) {
	value, err := usecase.transaction.DoInTx(ctx, func(ctx context.Context) (interface{}, error) {
		user, err := usecase.userRepository.FindByUniqueUser(ctx, userId, uuid)
		if err != nil {
			return nil, err
		}

		if user.UserKind == entity.Banned {
			c := &i18n.LocalizeConfig{MessageID: model.E0105}
			return nil, model.NewErrForbidden(model.E0105, usecase.localizer.MustLocalize(c))
		}

		if user.UpdateDevice(device, appVersion, platformNumber) {
			if err = usecase.userRepository.UpdateUser(ctx, user); err != nil {
				return nil, err
			}
		}

		accountToken, keyStr, ivStr, err := usecase.login(ctx, user)
		if err != nil {
			return nil, err
		}

		return model.NewSession(user, accountToken, keyStr, ivStr), nil
	})
	if err != nil {
		return nil, err
	}

	return value.(*model.Session), nil
}

func (usecase *sessionUsecase) login(ctx context.Context, user *entity.User) (string, string, string, error) {
	sessionID := uuid.New().String()
	accountToken, err := usecase.generateToken(user)
	if err != nil {
		c := &i18n.LocalizeConfig{MessageID: model.E9999}
		return "", "", "", model.NewErrUnprocessable(model.E9999, usecase.localizer.MustLocalize((c)))
	}

	key, iv, err := usecase.generateKeyAndIV()
	if err != nil {
		c := &i18n.LocalizeConfig{MessageID: model.E9999}
		return "", "", "", model.NewErrUnprocessable(model.E9999, usecase.localizer.MustLocalize((c)))
	}

	catData := append(key, iv...)
	cacheKey := formatter.CRYPTO_CACHE_KEY + sessionID
	err = usecase.cache.Set(ctx, cacheKey, catData, time.Hour*10)
	if err != nil {
		c := &i18n.LocalizeConfig{MessageID: model.E9999}
		return "", "", "", model.NewErrUnprocessable(model.E9999, usecase.localizer.MustLocalize((c)))
	}

	if !usecase.cfg.IsMultiDeviceAccess() {
		// ログイン別のSessionIDをキャッシュ
		sessionCat := []byte(sessionID)
		sessionCacheKey := formatter.CRYPTO_CACHE_KEY + user.UUID
		err = usecase.cache.Set(ctx, sessionCacheKey, sessionCat, time.Hour*10)
		if err != nil {
			c := &i18n.LocalizeConfig{MessageID: model.E9999}
			return "", "", "", model.NewErrUnprocessable(model.E9999, usecase.localizer.MustLocalize(c))
		}
	}

	keyStr := base64.StdEncoding.EncodeToString(key)
	ivStr := base64.StdEncoding.EncodeToString(iv)

	return accountToken, keyStr, ivStr, nil
}

func (usecase *sessionUsecase) generateToken(user *entity.User) (string, error) {
	claims := &middleware.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        strconv.FormatUint(uint64(user.ID), 10),
			ExpiresAt: jwt.NewNumericDate(flextime.Now().Add(10 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(flextime.Now()),
			NotBefore: jwt.NewNumericDate(flextime.Now()),
			Subject:   usecase.cfg.Jwt.Issuer,
			Issuer:    usecase.cfg.Jwt.Issuer,
			Audience:  []string{usecase.cfg.Jwt.Audience},
		},
		Uuid:        user.UUID,
		Name:        user.Name,
		InstalledAt: user.CreatedAt,
		CreatedAt:   flextime.Now(),
		UserKind:    user.UserKind,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := []byte(usecase.cfg.Jwt.Secret)

	return token.SignedString(key)
}

func (usecase *sessionUsecase) generateKeyAndIV() ([]byte, []byte, error) {
	key := make([]byte, formatter.KEY_SIZE)
	_, err := rand.Read(key)
	if err != nil {
		return nil, nil, err
	}

	iv := make([]byte, formatter.IV_SIZE)
	_, err = rand.Read(iv)
	if err != nil {
		return nil, nil, err
	}

	return key, iv, nil
}
