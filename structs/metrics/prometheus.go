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
	metrics  map[string]metricWrapper
	mu       sync.RWMutex
}

type metricWrapper struct {
	metric prometheus.Collector
	mType  metric_info.MetricType
}

// New 创建一个新的PrometheusMetrics实例
func New(namespace, subsystem string) *PrometheusMetrics {
	return &PrometheusMetrics{
		registry: prometheus.NewRegistry(),
		metrics:  make(map[string]metricWrapper),
		mu:       sync.RWMutex{},
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
		metric = prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: pm.namespace,
				Subsystem: pm.subsystem,
				Name:      info.Name,
				Help:      info.Help,
			},
		)
	case metric_info.Gauge:
		metric = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: pm.namespace,
				Subsystem: pm.subsystem,
				Name:      info.Name,
				Help:      info.Help,
			},
		)
	case metric_info.Histogram:
		metric = prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: pm.namespace,
				Subsystem: pm.subsystem,
				Name:      info.Name,
				Help:      info.Help,
				Buckets:   info.Buckets,
			},
		)
	case metric_info.Summary:
		metric = prometheus.NewSummary(
			prometheus.SummaryOpts{
				Namespace:  pm.namespace,
				Subsystem:  pm.subsystem,
				Name:       info.Name,
				Help:       info.Help,
				Objectives: info.Objectives,
			},
		)
	default:
		panic(fmt.Sprintf("[metrics] unknown metric type for [%s]", info.Name))
	}

	pm.registry.MustRegister(metric)
	pm.metrics[info.Name] = metricWrapper{metric: metric, mType: info.Type}
}

// GetMetric 通过名字获取指标
func (pm *PrometheusMetrics) getMetric(name string) (metricWrapper, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	metric, exists := pm.metrics[name]
	return metric, exists
}

// Report 上报数据
func (pm *PrometheusMetrics) Report(ctx context.Context, name string, labels map[string]string, value float64) {
	metricWrapper, exists := pm.getMetric(name)
	if !exists {
		logger.Warn(ctx, "metric not found", field.String("name", name))
		return
	}

	switch metricWrapper.mType {
	case metric_info.Counter:
		metricWrapper.metric.(prometheus.Counter).Add(value)
	case metric_info.Gauge:
		metricWrapper.metric.(prometheus.Gauge).Set(value)
	case metric_info.Histogram:
		metricWrapper.metric.(prometheus.Histogram).Observe(value)
	case metric_info.Summary:
		metricWrapper.metric.(prometheus.Summary).Observe(value)
	}
}

func (pm *PrometheusMetrics) GetRegistry() *prometheus.Registry {
	return pm.registry
}
