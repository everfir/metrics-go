package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/everfir/metrics-go"
)

// NetHTTPMetricsMiddleware 是针对 HTTP 协议的指标中间件
type NetHTTPMetricsMiddleware struct {
	*BaseMetricsMiddleware
}

// NewHTTPMetricsMiddleware 创建一个新的 HTTP 指标中间件
func HTTPMiddleware() *NetHTTPMetricsMiddleware {
	m := &NetHTTPMetricsMiddleware{
		BaseMetricsMiddleware: defaultBaseMiddleware.clone(),
	}
	m.Init()
	return m
}

// Middleware 返回一个适用于标准 net/http 的中间件函数
func (m *NetHTTPMetricsMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		labels := map[string]string{
			"method": r.URL.Path,
		}

		// 记录请求
		ctx := r.Context()
		metrics.Report(ctx, MetricRequestCnt, labels, 1)

		// 包装 ResponseWriter 以捕获状态码和响应大小
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// 调用下一个处理器
		next.ServeHTTP(rw, r)

		// 响应时间和状态码
		labels["status"] = strconv.Itoa(rw.statusCode)
		metrics.Report(ctx, MetricStatusCode, labels, 1)
		metrics.Report(ctx, MetricLatency, labels, float64(time.Since(start).Microseconds()))
	})
}

// responseWriter 是一个包装了 http.ResponseWriter 的结构体，用于捕获状态码和响应大小
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}
