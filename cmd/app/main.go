package main

import (
	"context"
	"html/template"
	"math"
	"net/http"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Daka-0424/my-go-server/config"
	controller "github.com/Daka-0424/my-go-server/pkg/controller/api"
	"github.com/Daka-0424/my-go-server/pkg/controller/route"
	"github.com/Daka-0424/my-go-server/pkg/controller/web/admin"
	"github.com/Daka-0424/my-go-server/pkg/controller/web/cheat"
	"github.com/Daka-0424/my-go-server/pkg/domain/entity"
	"github.com/Daka-0424/my-go-server/pkg/domain/service"
	"github.com/Daka-0424/my-go-server/pkg/infra"
	"github.com/Daka-0424/my-go-server/pkg/infra/appstore"
	"github.com/Daka-0424/my-go-server/pkg/infra/logger"
	"github.com/Daka-0424/my-go-server/pkg/infra/playstore"
	"github.com/Daka-0424/my-go-server/pkg/infra/repository"
	"github.com/Daka-0424/my-go-server/pkg/infra/util"
	"github.com/Daka-0424/my-go-server/pkg/usecase"
	"github.com/Masterminds/sprig/v3"
	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

// @title My Go Server API
// @version 1
// @description This is a sample server for My Go Server API.
// @securityDefinitions.apiKey Bearer
// @in header
// @name Authorization
func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	handler := initHandler(cfg)

	mysql := infra.NewMySQLConnector(cfg)
	redis := infra.NewRedisConnector(cfg)

	localizer := newLocalizer()

	migrate(mysql.DB)

	app := fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		fx.Supply(cfg, localizer, mysql.DB, redis.Client, handler, redis.RedSync),
		repository.Modules(),
		service.Modules(),
		usecase.Modules(),
		controller.Modules(),
		cheat.Modules(),
		admin.Modules(),
		appstore.Modules(),
		playstore.Modules(),
		logger.Modules(cfg),
		util.Modules(),
		fx.Invoke(
			lifecycle,
			route.Route,
			route.CheatRoute,
			route.AdminRoute,
			func(r *gin.Engine) {},
		),
	)

	app.Run()
}

func initHandler(cfg *config.Config) *gin.Engine {
	if !cfg.IsDevelopment() {
		gin.SetMode(gin.ReleaseMode)
	}

	securityCfg := secure.DefaultConfig()
	securityCfg.CustomFrameOptionsValue = "SAMEORIGIN"
	securityCfg.SSLRedirect = false
	securityCfg.ContentSecurityPolicy = ""

	handler := gin.New()

	handler.Use(secure.New(securityCfg))
	handler.Use(func(ctx *gin.Context) {
		ctx.Writer.Header().Set("X_Robots-Tag", "noindex, nofollow, nosnippet, noarchive")
	})
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	handler.SetFuncMap(customFunc())

	return handler
}

func customFunc() template.FuncMap {
	funcMap := sprig.FuncMap()
	funcMap["toHtmlDateTime"] = func(t time.Time) string { return t.Local().Format("2000-01-01T00:00:00") }
	funcMap["toUintMinute"] = func(d time.Duration) uint { return uint(d.Minutes()) }
	funcMap["toTruncation"] = func(c float64) uint { return uint(math.Round(c*100) / 100) }
	funcMap["toHtmlTime"] = func(t time.Time) string { return t.Format("00:00") }

	return funcMap
}

func newLocalizer() *i18n.Localizer {
	bundle := i18n.NewBundle(language.Japanese)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	_, err := bundle.LoadMessageFile("language/ja.toml")
	if err != nil {
		panic(err)
	}

	return i18n.NewLocalizer(bundle)
}

func lifecycle(lc fx.Lifecycle, cfg *config.Config, handler *gin.Engine) {
	srv := &http.Server{Addr: ":8080", Handler: handler}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go (func() {
				_ = srv.ListenAndServe()
			})()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
}

func migrate(db *gorm.DB) {
	err := db.AutoMigrate(entity.Entity()...)

	if err != nil {
		panic(err)
	}
}
