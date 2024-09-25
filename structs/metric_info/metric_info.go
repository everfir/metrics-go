package metric_info

// MetricType 定义了不同类型的指标
type MetricType int

const (
	Counter MetricType = iota
	Gauge
	Histogram
	Summary
)

// MetricInfo 封装了一个指标的所有配置信息
type MetricInfo struct {
	Type       MetricType          // 指标类型
	Name       string              // 指标名称
	Help       string              // 指标帮助信息
	Buckets    []float64           // 仅对Histogram有效
	Objectives map[float64]float64 // 仅对Summary有效
}
