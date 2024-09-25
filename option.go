package metrics

import (
	"time"

	"github.com/everfir/metrics-go/structs/config"
)

// Option 定义了一个函数类型，用于设置配置选项
type Option func(*config.MetricsConfig)

// WithCollectorMode 设置为 Collector 模式
func WithCollectorMode(port int) Option {
	return func(c *config.MetricsConfig) {
		c.ReportType = config.CollectorType
		c.Port = port
	}
}

// WithPushgatewayMode 设置为 Pushgateway 模式
func WithPushgatewayMode(pushAddr, jobName string, pushInterval time.Duration) Option {
	return func(c *config.MetricsConfig) {
		c.ReportType = config.PushgatewayType
		c.PushAddr = pushAddr
		c.JobName = jobName
		c.PushInterval = pushInterval
	}
}

// WithNamespace 设置 namespace
func WithNamespace(namespace string) Option {
	return func(c *config.MetricsConfig) {
		c.Namespace = namespace
	}
}

// WithSubsystem 设置 subsystem
func WithSubsystem(subsystem string) Option {
	return func(c *config.MetricsConfig) {
		c.Subsystem = subsystem
	}
}
