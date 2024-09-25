package reporter

import (
	"context"
	"fmt"
	"time"

	"github.com/everfir/metrics-go/structs/metric_info"
	"github.com/everfir/metrics-go/structs/metrics"
	"github.com/prometheus/client_golang/prometheus/push"
)

type PushgatewayReporter struct {
	metrics   *metrics.PrometheusMetrics
	pusher    *push.Pusher
	pushAddr  string
	jobName   string
	pushTimer *time.Ticker
}

func NewPushgatewayReporter(namespace, subsystem, pushAddr, jobName string, pushInterval time.Duration) *PushgatewayReporter {
	m := metrics.New(namespace, subsystem)
	pusher := push.New(pushAddr, jobName).Gatherer(m.GetRegistry())

	reporter := &PushgatewayReporter{
		metrics:   m,
		pusher:    pusher,
		pushAddr:  pushAddr,
		jobName:   jobName,
		pushTimer: time.NewTicker(pushInterval),
	}

	go reporter.startPushing()

	return reporter
}

func (p *PushgatewayReporter) startPushing() {
	for range p.pushTimer.C {
		if err := p.pusher.Push(); err != nil {
			fmt.Printf("Could not push to Pushgateway: %v\n", err)
		}
	}
}

func (p *PushgatewayReporter) Register(info metric_info.MetricInfo) {
	p.metrics.Register(info)
}

func (p *PushgatewayReporter) Report(ctx context.Context, name string, labels map[string]string, value float64) {
	p.metrics.Report(ctx, name, labels, value)
}

func (p *PushgatewayReporter) Close(ctx context.Context) error {
	p.pushTimer.Stop()
	return nil
}
