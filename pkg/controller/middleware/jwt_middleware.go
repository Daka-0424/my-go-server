package middleware

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/controller/formatter"
	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Claims struct {
	jwt.RegisteredClaims
	Uuid        string    `json:"uuid"`
	Name        string    `json:"name"`
	InstalledAt time.Time `json:"installed_at"`
	CreatedAt   time.Time `json:"created_at"`
	UserKind    uint      `json:"user_kind"`
}

func (c Claims) GetUserId() (uint, error) {
	id, err := strconv.ParseUint(c.ID, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func (c Claims) IsSuperUser() bool {
	return c.UserKind == entity.SuperUser
}

func JwtMiddleware(cfg *config.Config, localizer *i18n.Localizer, cache repository.ICache) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tknstr, err := bearerToken(ctx)
		if err != nil {
			returnErrorWithAbort(ctx, localizer)
			return
		}

		claims := &Claims{}

		tkn, err := jwt.ParseWithClaims(tknstr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.Jwt.Secret), nil
		})
		if err != nil {
			returnErrorWithAbort(ctx, localizer)
			return
		}

		if !tkn.Valid {
			returnErrorWithAbort(ctx, localizer)
			return
		}

		ctx.Set("claims", claims)

		data, ok, err := cache.Get(ctx, formatter.CRYPTO_CACHE_KEY+claims.Uuid)
		if err != nil {
			returnErrorWithAbort(ctx, localizer)
			return
		}

		if ok {
			ctx.Set("cryptoKey", data[:formatter.KEY_SIZE])
			ctx.Set("cryptoIv", data[formatter.KEY_SIZE:])
		}
	}
}

func bearerToken(ctx *gin.Context) (string, error) {
	auth := ctx.Request.Header.Get("Authorization")
	if auth == "" {
		return "", errors.New("token not found")
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == "" {
		return "", errors.New("token not found")
	}
	return token, nil
}

func returnErrorWithAbort(ctx *gin.Context, localizer *i18n.Localizer) {
	c := &i18n.LocalizeConfig{MessageID: model.E1001}
	appErr := model.NewErrUnauthorized(model.E1001, localizer.MustLocalize(c))
	formatter.Respond(ctx, appErr.StatusCode, gin.H{"error": appErr})
	ctx.Abort()
}
