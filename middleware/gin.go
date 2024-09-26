package middleware

import (
	"strconv"
	"time"

	"github.com/everfir/metrics-go"
	"github.com/gin-gonic/gin"
)

// GinMetricsMiddleware 是针对 Gin 框架的指标中间件
type GinMetricsMiddleware struct {
	*BaseMetricsMiddleware
}

// GinMiddleware 创建一个新的 Gin 指标中间件
func GinMiddleware() *GinMetricsMiddleware {
	m := &GinMetricsMiddleware{
		BaseMetricsMiddleware: defaultBaseMiddleware.clone(),
	}
	return m
}

// Middleware 返回一个适用于 Gin 框架的中间件函数
func (m *GinMetricsMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		labels := map[string]string{
			"method": c.Request.URL.Path,
		}

		// 记录请求
		metrics.Report(c, MetricRequestCnt, labels, 1)

		// 调用下一个处理器
		c.Next()

		// 记录指标
		labels["status"] = strconv.Itoa(c.Writer.Status())
		metrics.Report(c, MetricStatusCode, labels, 1)
		metrics.Report(c, MetricLatency, labels, float64(time.Since(start).Microseconds()))
	}
}
