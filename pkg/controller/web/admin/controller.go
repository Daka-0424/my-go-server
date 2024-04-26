package admin

import (
	"encoding/json"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/repository"
	"github.com/Daka-0424/my-go-server/pkg/usecase/model"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const (
	PageSize = 50
)

type adminControllerBase struct {
	cfg       *config.Config
	cache     repository.Cache
	localizer *i18n.Localizer
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
		cfg := &i18n.LocalizeConfig{MessageID: model.E2004}
		return nil, model.NewErrUnprocessable(model.E2004, ctl.localizer.MustLocalize(cfg))
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
