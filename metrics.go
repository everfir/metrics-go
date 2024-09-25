package metrics

import (
	"context"
	"fmt"
	"sync"

	"github.com/everfir/metrics-go/structs/config"
	"github.com/everfir/metrics-go/structs/metric_info"
	"github.com/everfir/metrics-go/structs/reporter"
)

var (
	r    reporter.MetricsReporter
	once sync.Once
)

// Init 初始化 metrics 系统
func Init(opts ...Option) error {
	cfg := &config.MetricsConfig{
		ReportType: config.CollectorType, // 默认使用 Collector 模式
		Port:       8080,                 // 默认端口
	}

	for _, opt := range opts {
		opt(cfg)
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	var err error
	once.Do(func() {
		switch cfg.ReportType {
		case config.CollectorType:
			r = reporter.NewCollectorReporter(cfg.Namespace, cfg.Subsystem, cfg.Port)
		case config.PushgatewayType:
			r = reporter.NewPushgatewayReporter(cfg.Namespace, cfg.Subsystem, cfg.PushAddr, cfg.JobName, cfg.PushInterval)
		default:
			err = fmt.Errorf("invalid report type")
			return
		}
	})

	return err
}

// Close 优雅地关闭metrics系统
func Close(ctx context.Context) error {
	if r != nil {
		return r.Close(ctx)
	}
	return nil
}

// Register 允许用户注册新的指标
func Register(info metric_info.MetricInfo) {
	if r == nil {
		panic("[metrics] metrics not initialized, call Init() first")
	}
	r.Register(info)
}

// Report 允许用户上报数据
func Report(ctx context.Context, name metric_info.MetricName, labels map[string]string, value float64) {
	if r == nil {
		panic("metrics not initialized, call Init() first")
	}
	r.Report(ctx, name, labels, value)
}
