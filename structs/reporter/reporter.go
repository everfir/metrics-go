package reporter

import (
	"context"

	"github.com/everfir/metrics-go/structs/metric_info"
)

// MetricsReporter 定义了指标上报的接口
type MetricsReporter interface {
	Register(info metric_info.MetricInfo)
	Report(ctx context.Context, name metric_info.MetricName, labels map[string]string, value float64)
	Close(ctx context.Context) error
}
