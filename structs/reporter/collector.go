package reporter

import (
	"context"
	"fmt"
	"net/http"

	"github.com/everfir/metrics-go/structs/metric_info"
	"github.com/everfir/metrics-go/structs/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type CollectorReporter struct {
	metrics *metrics.PrometheusMetrics
	server  *http.Server
}

func NewCollectorReporter(namespace, subsystem string, port int) *CollectorReporter {
	m := metrics.New(namespace, subsystem)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(m.GetRegistry(), promhttp.HandlerOpts{}))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Sprintf("Failed to start metrics server: %v", err))
		}
	}()

	return &CollectorReporter{
		metrics: m,
		server:  srv,
	}
}

func (c *CollectorReporter) Register(info metric_info.MetricInfo) {
	c.metrics.Register(info)
}

func (c *CollectorReporter) Report(ctx context.Context, name metric_info.MetricName, labels map[string]string, value float64) {
	c.metrics.Report(ctx, name, labels, value)
}

func (c *CollectorReporter) Close(ctx context.Context) error {
	return c.server.Shutdown(ctx)
}
