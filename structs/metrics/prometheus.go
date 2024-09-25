package metrics

import (
	"context"
	"fmt"
	"sync"

	"github.com/everfir/logger-go"
	"github.com/everfir/logger-go/structs/field"
	"github.com/everfir/metrics-go/structs/metric_info"
	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusMetrics 实现了MetricsInterface，使用Prometheus Go SDK
type PrometheusMetrics struct {
	namespace string
	subsystem string

	registry *prometheus.Registry
	metrics  map[metric_info.MetricName]metricWrapper
	mu       sync.RWMutex
}

type metricWrapper struct {
	metric prometheus.Collector
	info   metric_info.MetricInfo
}

// New 创建一个新的PrometheusMetrics实例
func New(namespace, subsystem string) *PrometheusMetrics {
	return &PrometheusMetrics{
		namespace: namespace,
		subsystem: subsystem,
		registry:  prometheus.NewRegistry(),
		metrics:   make(map[metric_info.MetricName]metricWrapper),
		mu:        sync.RWMutex{},
	}
}

// Register 根据MetricInfo自动注册指标
func (pm *PrometheusMetrics) Register(info metric_info.MetricInfo) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// 不能重复注册
	if _, exists := pm.metrics[info.Name]; exists {
		panic(fmt.Sprintf("[metrics] metric [%s] already registered with a different type", info.Name))
	}

	var metric prometheus.Collector
	switch info.Type {
	case metric_info.Counter:
		metric = prometheus.V2.NewCounterVec(
			prometheus.CounterVecOpts{
				CounterOpts: prometheus.CounterOpts{
					Namespace: pm.namespace,
					Subsystem: pm.subsystem,
					Name:      info.Name.String(),
					Help:      info.Help,
				},
				VariableLabels: info.ToConstrainableLabels(),
			},
		)
	case metric_info.Gauge:
		metric = prometheus.V2.NewGaugeVec(
			prometheus.GaugeVecOpts{
				GaugeOpts: prometheus.GaugeOpts{
					Namespace: pm.namespace,
					Subsystem: pm.subsystem,
					Name:      info.Name.String(),
					Help:      info.Help,
				},
				VariableLabels: info.ToConstrainableLabels(),
			},
		)
	case metric_info.Histogram:
		metric = prometheus.V2.NewHistogramVec(
			prometheus.HistogramVecOpts{
				HistogramOpts: prometheus.HistogramOpts{
					Namespace: pm.namespace,
					Subsystem: pm.subsystem,
					Name:      info.Name.String(),
					Help:      info.Help,
					Buckets:   info.Buckets,
				},
				VariableLabels: info.ToConstrainableLabels(),
			},
		)
	case metric_info.Summary:
		metric = prometheus.V2.NewSummaryVec(
			prometheus.SummaryVecOpts{
				SummaryOpts: prometheus.SummaryOpts{
					Namespace:  pm.namespace,
					Subsystem:  pm.subsystem,
					Name:       info.Name.String(),
					Help:       info.Help,
					Objectives: info.Objectives,
				},
				VariableLabels: info.ToConstrainableLabels(),
			},
		)
	default:
		panic(fmt.Sprintf("[metrics] unknown metric type for [%s]", info.Name))
	}

	pm.registry.MustRegister(metric)
	pm.metrics[info.Name] = metricWrapper{metric: metric, info: info}
}

// GetMetric 通过名字获取指标
func (pm *PrometheusMetrics) getMetric(name metric_info.MetricName) (metricWrapper, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	metric, exists := pm.metrics[name]
	return metric, exists
}

// Report 上报数据
func (pm *PrometheusMetrics) Report(ctx context.Context, name metric_info.MetricName, labels map[string]string, value float64) {
	// 获取指标包装器
	metricWrapper, exists := pm.getMetric(name)
	if !exists {
		// 如果指标不存在，记录警告日志并返回
		logger.Warn(ctx, "metric not found", field.String("name", name.String()))
		return
	}

	// 创建标签映射
	var mapping = map[string]string{}
	// 处理预定义的标签处理函数
	for k, v := range metricWrapper.info.LabelHandler {
		mapping[k] = v(ctx)
	}
	// 合并用户提供的标签
	for k, v := range labels {
		mapping[k] = v
	}

	// 根据指标类型进行不同的处理
	switch metricWrapper.info.Type {
	case metric_info.Counter:
		// 对于计数器类型，增加指定的值
		metricWrapper.metric.(*prometheus.CounterVec).With(mapping).Add(value)
	case metric_info.Gauge:
		// 对于仪表类型，设置指定的值
		metricWrapper.metric.(*prometheus.GaugeVec).With(mapping).Set(value)
	case metric_info.Histogram:
		// 对于直方图类型，观察指定的值
		metricWrapper.metric.(*prometheus.HistogramVec).With(mapping).Observe(value)
	case metric_info.Summary:
		// 对于摘要类型，观察指定的值
		metricWrapper.metric.(*prometheus.SummaryVec).With(mapping).Observe(value)
	default:
		// 对于未知类型，记录错误日志
		logger.Error(ctx, "未知的指标类型",
			field.String("name", name.String()),
			field.String("type", string(metricWrapper.info.Type)))
	}
}

func (pm *PrometheusMetrics) GetRegistry() *prometheus.Registry {
	return pm.registry
}
