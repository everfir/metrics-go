package metric_info

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

// MetricType 定义了不同类型的指标
type MetricType int

const (
	Counter MetricType = iota
	Gauge
	Histogram
	Summary
)

func (m MetricType) String() string {
	switch m {
	case Counter:
		return "counter"
	case Gauge:
		return "gauge"
	case Histogram:
		return "histogram"
	case Summary:
		return "summary"
	default:
		return "unknown"
	}
}

type MetricName string

func (m MetricName) String() string {
	return string(m)
}

func NewMetricName(name string) MetricName {
	return MetricName(name)
}

// MetricInfo 封装了一个指标的所有配置信息
type MetricInfo struct {
	Type       MetricType          // 指标类型
	Name       MetricName          // 指标名称
	Help       string              // 指标帮助信息
	Buckets    []float64           // 仅对Histogram有效
	Objectives map[float64]float64 // 仅对Summary有效

	// 标签「必选」
	Labels       []string
	LabelHandler map[string]LabelHandler
}

// LabelHandler 定义为一个函数类型，接收上下文，返回标签值
type LabelHandler func(ctx context.Context) string

// ToConstrainableLabels 将 LabelHandler 转换为 Prometheus 的 ConstrainableLabels
func (mi *MetricInfo) ToConstrainableLabels() (ret prometheus.ConstrainedLabels) {
	for _, name := range mi.Labels {
		ret = append(ret, prometheus.ConstrainedLabel{
			Name: name,
		})
	}
	return ret
}
