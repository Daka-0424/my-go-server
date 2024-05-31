package logger

import "context"

type IKpiLoggerFactory interface {
	Create(ctx context.Context) (IKpi, error)
}
