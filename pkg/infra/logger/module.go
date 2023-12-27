package logger

import (
	"github.com/Daka-0424/my-go-server/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Modules(cfg *config.Config) fx.Option {
	zapLoggerFactory := zap.NewProduction
	if cfg.IsDevelopment() {
		zapLoggerFactory = zap.NewDevelopment
	}

	return fx.Module("logger",
		fx.Provide(
			zapLoggerFactory,
		),
	)
}
