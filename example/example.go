package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/everfir/metrics-go"
	"github.com/everfir/metrics-go/middleware"
	"github.com/everfir/metrics-go/structs/metric_info"
	"github.com/gin-gonic/gin"
)

func main() {
	// 定义命令行参数
	mode := flag.String("mode", "collector", "Exporter mode: collector or pushgateway")
	port := flag.Int("port", 10086, "Port for collector mode")
	pushAddr := flag.String("push-addr", "http://localhost:9091", "Pushgateway address for pushgateway mode")
	pushInterval := flag.Duration("push-interval", 5*time.Second, "Push interval for pushgateway mode")
	serverType := flag.String("server", "http", "Server type: http or gin")
	flag.Parse()

	// 准备通用的选项
	opts := []metrics.Option{
		metrics.WithNamespace("everfir"),
		metrics.WithSubsystem("metrics_example"),
	}

	// 根据模式添加特定的选项
	switch *mode {
	case "collector":
		opts = append(opts, metrics.WithCollectorMode(*port))
	case "pushgateway":
		opts = append(opts, metrics.WithPushgatewayMode(*pushAddr, "example_job", *pushInterval))
	default:
		fmt.Println("Invalid mode. Use 'collector' or 'pushgateway'.")
		os.Exit(1)
	}

	// 初始化 metrics 系统
	err := metrics.Init(opts...)
	if err != nil {
		fmt.Printf("Failed to initialize metrics: %v\n", err)
		os.Exit(1)
	}

	// 根据模式输出信息
	switch *mode {
	case "collector":
		fmt.Printf("Collector mode: Metrics server started on port %d\n", *port)
	case "pushgateway":
		fmt.Printf("Pushgateway mode: Pushing to %s every %v\n", *pushAddr, *pushInterval)
	}

	// 注册指标
	metrics.Register(metric_info.MetricInfo{
		Type:   metric_info.Counter,
		Name:   metric_info.MetricName("example_counter"),
		Help:   "An example counter",
		Labels: []string{"label"},
	})

	metrics.Register(metric_info.MetricInfo{
		Type:   metric_info.Gauge,
		Name:   metric_info.MetricName("example_gauge"),
		Help:   "An example gauge",
		Labels: []string{"label"},
	})

	metrics.Register(metric_info.MetricInfo{
		Type:    metric_info.Histogram,
		Name:    metric_info.MetricName("example_histogram"),
		Help:    "An example histogram",
		Buckets: []float64{1, 5, 10, 50, 100},
		Labels:  []string{"label"},
	})

	// 创建一个通道来接收终止信号
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器
	switch *serverType {
	case "http":
		go startHTTPServer()
	case "gin":
		go startGinServer()
	default:
		fmt.Println("Invalid server type. Use 'http' or 'gin'.")
		os.Exit(1)
	}

	fmt.Println("Application is running. Press Ctrl+C to exit.")

	// 等待终止信号
	<-done
	fmt.Println("Received termination signal. Shutting down...")

	// 优雅地关闭 metrics 系统
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := metrics.Close(ctx); err != nil {
		fmt.Printf("Error shutting down metrics system: %v\n", err)
	}

	fmt.Println("Application has been shut down.")
}

func startHTTPServer() {
	httpMiddleware := middleware.HTTPMiddleware()
	http.Handle("/", httpMiddleware.Middleware(http.HandlerFunc(handler)))
	http.ListenAndServe(":8080", nil)
}

func startGinServer() {
	r := gin.Default()
	ginMiddleware := middleware.GinMiddleware()
	r.Use(ginMiddleware.Middleware())
	r.GET("/测试gin", ginHandler)
	r.Run(":8080")
}

func handler(w http.ResponseWriter, r *http.Request) {
	// 模拟一些数据上报
	counterValue := 1.0
	gaugeValue := rand.Float64() * 100
	histogramValue := rand.Float64() * 200

	metrics.Report(r.Context(), metric_info.MetricName("example_counter"), map[string]string{"label": "value"}, counterValue)
	metrics.Report(r.Context(), metric_info.MetricName("example_gauge"), map[string]string{"label": "value"}, gaugeValue)
	metrics.Report(r.Context(), metric_info.MetricName("example_histogram"), map[string]string{"label": "value"}, histogramValue)

	fmt.Fprintf(w, "Reported metrics - Counter: %.2f, Gauge: %.2f, Histogram: %.2f\n", counterValue, gaugeValue, histogramValue)
}

func ginHandler(c *gin.Context) {
	// 模拟一些数据上报
	counterValue := 1.0
	gaugeValue := rand.Float64() * 100
	histogramValue := rand.Float64() * 200

	metrics.Report(c.Request.Context(), metric_info.MetricName("example_counter"), map[string]string{"label": "value"}, counterValue)
	metrics.Report(c.Request.Context(), metric_info.MetricName("example_gauge"), map[string]string{"label": "value"}, gaugeValue)
	metrics.Report(c.Request.Context(), metric_info.MetricName("example_histogram"), map[string]string{"label": "value"}, histogramValue)

	c.String(200, "Reported metrics - Counter: %.2f, Gauge: %.2f, Histogram: %.2f\n", counterValue, gaugeValue, histogramValue)
}
