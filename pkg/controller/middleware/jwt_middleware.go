package middleware

import (
	"strconv"
	"time"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/language"
	"github.com/Daka-0424/my-go-server/pkg/controller/formatter"
	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	SessionID   string    `json:"session_id"`
	Uuid        string    `json:"uuid"`
	Name        string    `json:"name"`
	InstalledAt time.Time `json:"installed_at"`
	CreatedAt   time.Time `json:"created_at"`
	UserKind    uint      `json:"user_kind"`
}

func (c Claims) GetUserId(ctx *gin.Context, localizer *language.Localizer) (uint, *response.AppError) {
	id, err := strconv.ParseUint(c.ID, 10, 64)
	if err != nil {
		lang := ctx.Request.Header.Get("Accept-Language")
		return 0, response.NewErrInternalServerError(model.E0101, localizer.MustLocalize(model.E0101, lang, nil))
	}
	return uint(id), nil
}

func (c Claims) IsSuperUser() bool {
	return c.UserKind == entity.SuperUser
}

func JwtMiddleware(cfg *config.Config, localizer *language.Localizer, cache repository.ICache) gin.HandlerFunc {
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

		data, ok, err := cache.Get(ctx, formatter.CRYPTO_CACHE_KEY+claims.SessionID)
		if err != nil || !ok {
			returnErrorWithAbort(ctx, localizer)
			return
		}

		ctx.Set("cryptoKey", data[:formatter.KEY_SIZE])
		ctx.Set("cryptoIv", data[formatter.KEY_SIZE:])

		if !cfg.IsMultiDeviceAccess() {
			session, ok, err := cache.Get(ctx, formatter.CRYPTO_CACHE_KEY+claims.Uuid)
			if err != nil || !ok {
				returnErrorWithAbort(ctx, localizer)
				return
			}

			if claims.SessionID != string(session) {
				returnErrorWithAbort(ctx, localizer)
				return
			}
		}
	}
}
