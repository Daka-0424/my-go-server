package logger

import (
	"context"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/domain/logger"
)

type kpiLoggerFactory struct {
	cfg *config.Config
}

func NewKpiLoggerFactory(cfg *config.Config) logger.IKpiLoggerFactory {
	return &kpiLoggerFactory{
		cfg: cfg,
	}
}

func (factory *kpiLoggerFactory) Create(ctx context.Context) (logger.IKpi, error) {
	if factory.cfg.Kpi.ProjectID != "" {
		return NewKpiLogger(ctx, factory.cfg)
	}

	return NewLocalKpiLogger(ctx, factory.cfg)
}
