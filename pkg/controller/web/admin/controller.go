package admin

import (
	"encoding/json"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/language"
	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model/response"
	"github.com/gin-gonic/gin"
)

const (
	PageSize = 50
)

type adminControllerBase struct {
	cfg       *config.Config
	cache     repository.ICache
	localizer *language.Localizer
}

func (ctl *adminControllerBase) baseObj(ctx *gin.Context) gin.H {
	admin, _ := ctl.getSession(ctx)
	return gin.H{
		"admin": admin,
		"env":   ctl.cfg.Settings.Environment,
	}
}

func (ctl *adminControllerBase) getSession(ctx *gin.Context) (*entity.Admin, error) {
	key := ctl.cfg.Cookie.Key
	redisKey, err := ctx.Cookie(key)
	if err != nil {
		return nil, err
	}

	redisValue, ok, err := ctl.cache.Get(ctx, redisKey)
	if err != nil {
		return nil, err
	}

	if !ok {
		lang := ctx.Request.Header.Get("Accept-Language")
		return nil, response.NewErrUnprocessable(model.E2004, ctl.localizer.MustLocalize(model.E2004, lang, nil))
	}

	admin := &entity.Admin{}
	if err := json.Unmarshal(redisValue, admin); err != nil {
		return nil, err
	}

	return admin, nil
}

func flashMessage(status, message string) gin.H {
	return gin.H{
		"status":  status,
		"message": message,
	}
}
