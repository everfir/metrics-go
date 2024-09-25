package config

import (
	"fmt"
	"time"
)

// ReportType 定义了指标上报的类型
type ReportType int

const (
	CollectorType ReportType = iota
	PushgatewayType
)

// MetricsConfig 包含所有配置选项
type MetricsConfig struct {
	ReportType   ReportType
	Namespace    string
	Subsystem    string
	Port         int
	PushAddr     string
	JobName      string
	PushInterval time.Duration
}

// Validate 验证配置的有效性
func (c *MetricsConfig) Validate() error {
	if c.Namespace == "" {
		return fmt.Errorf("namespace cannot be empty")
	}
	if c.Subsystem == "" {
		return fmt.Errorf("subsystem cannot be empty")
	}
	switch c.ReportType {
	case CollectorType:
		if c.Port <= 0 || c.Port > 65535 {
			return fmt.Errorf("invalid port number: %d", c.Port)
		}
	case PushgatewayType:
		if c.PushAddr == "" {
			return fmt.Errorf("pushAddr cannot be empty for Pushgateway mode")
		}
		if c.JobName == "" {
			return fmt.Errorf("jobName cannot be empty for Pushgateway mode")
		}
		if c.PushInterval <= 0 {
			return fmt.Errorf("pushInterval must be positive")
		}
	default:
		return fmt.Errorf("invalid report type: %v", c.ReportType)
	}
	return nil
}
