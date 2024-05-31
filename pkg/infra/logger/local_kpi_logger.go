package logger

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Daka-0424/my-go-server/config"
	"github.com/Daka-0424/my-go-server/pkg/domain/logger"
)

type localKpiLogger struct {
}

func NewLocalKpiLogger(ctx context.Context, cfg *config.Config) (logger.IKpi, error) {
	kpiLogger := &localKpiLogger{}

	return kpiLogger, nil
}

func (logger *localKpiLogger) LogEvent(_ logger.KpiLogName, data map[string]interface{}) {
	data["datetime"] = time.Now()
	logger.log(data)
}

func (logger *localKpiLogger) log(payload interface{}) {
	json, err := json.Marshal(payload)
	if err != nil {
		println(err.Error())
		return
	}

	println(string(json))
}

func (*localKpiLogger) Flush() error {
	return nil
}
