package logger

import (
	"context"
	"time"

	"cloud.google.com/go/logging"
	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/domain/logger"
)

type kpiLogger struct {
	client *logging.Client
}

func NewKpiLogger(ctx context.Context, cfg *config.Config) (logger.IKpi, error) {
	lg, err := logging.NewClient(ctx, cfg.Kpi.ProjectID)
	if err != nil {
		return nil, err
	}

	kpiLogger := &kpiLogger{
		client: lg,
	}

	return kpiLogger, nil
}

func (logger *kpiLogger) LogEvent(logName logger.KpiLogName, data map[string]interface{}) {
	data["datetime"] = time.Now()
	logger.log(logName, data)
}

func (logger *kpiLogger) log(logName logger.KpiLogName, payload interface{}) {
	buffer := logger.client.Logger(string(logName))
	buffer.Log(logging.Entry{Payload: payload})
}

func (logger *kpiLogger) Flush() error {
	return logger.client.Close()
}
