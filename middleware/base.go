package middleware

import (
	"context"

	"github.com/everfir/metrics-go"
	"github.com/everfir/metrics-go/structs/metric_info"
)

const (
	MetricLatency    metric_info.MetricName = metric_info.MetricName("latency")
	MetricStatusCode metric_info.MetricName = metric_info.MetricName("status_code")
	MetricRequestCnt metric_info.MetricName = metric_info.MetricName("req_cnt")
)

var (
	defaultBaseMiddleware = NewBaseMetricsMiddleware()
)

// BaseMetricsMiddleware 包含所有协议共用的功能
type BaseMetricsMiddleware struct {
	buildinMetrics map[metric_info.MetricName]*metric_info.MetricInfo
}

// NewBaseMetricsMiddleware 创建一个新的 BaseMetricsMiddleware
func NewBaseMetricsMiddleware() (ret *BaseMetricsMiddleware) {
	ret = &BaseMetricsMiddleware{
		buildinMetrics: make(map[metric_info.MetricName]*metric_info.MetricInfo),
	}

	ret.buildinMetrics[MetricRequestCnt] = &metric_info.MetricInfo{
		Type:         metric_info.Counter,
		Name:         MetricRequestCnt,
		Help:         "请求总数",
		Labels:       []string{"method"},
		LabelHandler: map[string]metric_info.LabelHandler{},
	}

	ret.buildinMetrics[MetricLatency] = &metric_info.MetricInfo{
		Type:         metric_info.Histogram,
		Name:         MetricLatency,
		Help:         "请求时延",
		Labels:       []string{"method", "status"},
		LabelHandler: map[string]metric_info.LabelHandler{},
	}

	ret.buildinMetrics[MetricStatusCode] = &metric_info.MetricInfo{
		Type:         metric_info.Counter,
		Name:         MetricStatusCode,
		Help:         "响应状态码",
		Labels:       []string{"method", "status"},
		LabelHandler: map[string]metric_info.LabelHandler{},
	}

	return
}

// WithLabel 添加一个标签处理器
func (b *BaseMetricsMiddleware) WithLabelHandler(label string, handler metric_info.LabelHandler) {
	for _, info := range b.buildinMetrics {
		if _, ok := info.LabelHandler[label]; !ok {
			info.Labels = append(info.Labels, label)
		}
		info.LabelHandler[label] = handler
	}
}

func (b *BaseMetricsMiddleware) WithMetric(info *metric_info.MetricInfo) {
	b.buildinMetrics[info.Name] = info
}

// UpdateMetric 更新一个指标， 如果指标不存在，则注册一个新指标，仅能在初始化时调用
func (b *BaseMetricsMiddleware) UpdateMetric(name metric_info.MetricName, info *metric_info.MetricName) {
	if info, ok := b.buildinMetrics[name]; ok {
		b.buildinMetrics[info.Name] = info
		delete(b.buildinMetrics, name)
	}
}

func (b *BaseMetricsMiddleware) clone() *BaseMetricsMiddleware {
	clone := &BaseMetricsMiddleware{
		buildinMetrics: make(map[metric_info.MetricName]*metric_info.MetricInfo),
	}
	for name, info := range b.buildinMetrics {
		clone.buildinMetrics[name] = info
	}
	return clone
}

func (b *BaseMetricsMiddleware) Init(ctx context.Context) {
	for _, info := range b.buildinMetrics {
		metrics.Register(ctx, *info)
	}
}
