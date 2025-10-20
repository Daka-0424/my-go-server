package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strconv"
	"time"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/language"
	"github.com/Daka-0424/my-go-server/pkg/controller/formatter"
	"github.com/Daka-0424/my-go-server/pkg/controller/middleware"
	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/logger"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model/request"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model/response"
	"github.com/Songmu/flextime"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type ISession interface {
	CreateSession(ctx context.Context, device, appVersion string, platformNumber uint, req request.Session) (*response.Session, error)
}

type sessionUsecase struct {
	cfg              *config.Config
	localizer        *language.Localizer
	cache            repository.ICache
	transaction      repository.ITransaction
	userRepository   repository.IUser
	kpiLoggerFactory logger.IKpiLoggerFactory
}

func NewSessionUsecase(
	cfg *config.Config,
	lc *language.Localizer,
	cache repository.ICache,
	transaction repository.ITransaction,
	userRepository repository.IUser,
	kpiLoggerFactory logger.IKpiLoggerFactory,
) ISession {
	return &sessionUsecase{
		cfg:              cfg,
		localizer:        lc,
		cache:            cache,
		transaction:      transaction,
		userRepository:   userRepository,
		kpiLoggerFactory: kpiLoggerFactory,
	}
}

func (usecase *sessionUsecase) CreateSession(ctx context.Context, device, appVersion string, platformNumber uint, req request.Session) (*response.Session, error) {
	value, err := usecase.transaction.DoInTx(ctx, func(ctx context.Context) (interface{}, error) {
		user, err := usecase.userRepository.FindByUniqueUser(ctx, req.UserId, req.Uuid)
		if err != nil {
			return nil, err
		}

		if user.UserKind == entity.Banned {
			return nil, response.NewErrForbidden(model.E0105, usecase.localizer.MustLocalize(model.E0105, language.LanguageJapanese, nil))
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

		return response.NewSession(user, accountToken, keyStr, ivStr), nil
	})
	if err != nil {
		return nil, err
	}

	return value.(*response.Session), nil
}

func (usecase *sessionUsecase) login(ctx context.Context, user *entity.User) (string, string, string, error) {
	kpiLogger, err := usecase.kpiLoggerFactory.Create(ctx)
	if err != nil {
		return "", "", "", response.NewErrBadRequest(model.E9999, usecase.localizer.MustLocalize(model.E9999, language.LanguageJapanese, nil))
	}

	sessionID := uuid.New().String()
	accountToken, err := usecase.generateToken(user, sessionID)
	if err != nil {
		return "", "", "", response.NewErrUnprocessable(model.E9999, usecase.localizer.MustLocalize(model.E9999, language.LanguageJapanese, nil))
	}

	key, iv, err := usecase.generateKeyAndIV()
	if err != nil {
		return "", "", "", response.NewErrUnprocessable(model.E9999, usecase.localizer.MustLocalize(model.E9999, language.LanguageJapanese, nil))
	}

	catData := append(key, iv...)
	cacheKey := formatter.CRYPTO_CACHE_KEY + sessionID
	err = usecase.cache.Set(ctx, cacheKey, catData, time.Hour*10)
	if err != nil {
		return "", "", "", response.NewErrUnprocessable(model.E9999, usecase.localizer.MustLocalize(model.E9999, language.LanguageJapanese, nil))
	}

	if !usecase.cfg.IsMultiDeviceAccess() {
		// ログイン別のSessionIDをキャッシュ
		sessionCat := []byte(sessionID)
		sessionCacheKey := formatter.CRYPTO_CACHE_KEY + user.UUID
		err = usecase.cache.Set(ctx, sessionCacheKey, sessionCat, time.Hour*10)
		if err != nil {
			return "", "", "", response.NewErrUnprocessable(model.E9999, usecase.localizer.MustLocalize(model.E9999, language.LanguageJapanese, nil))
		}
	}

	keyStr := base64.StdEncoding.EncodeToString(key)
	ivStr := base64.StdEncoding.EncodeToString(iv)

	kpiDate := map[string]interface{}{
		"user_id": user.ID,
	}
	kpiLogger.LogEvent(logger.KpiLogLogin, kpiDate)
	kpiLogger.Flush()

	return accountToken, keyStr, ivStr, nil
}

func (usecase *sessionUsecase) generateToken(user *entity.User, sessionID string) (string, error) {
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
		SessionID:   sessionID,
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
