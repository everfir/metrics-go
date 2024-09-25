# Everfir Metrics-Go

Everfir Metrics-Go 是一个灵活的 Go 语言指标收集和报告库，支持 Prometheus 的 Collector 和 Pushgateway 两种模式。该库提供了简单易用的 API，使得在 Go 应用程序中集成指标收集变得简单高效。

## 特性

- 支持 Prometheus Collector 模式
- 支持 Prometheus Pushgateway 模式
- 简单易用的 API
- 支持 Counter、Gauge、Histogram 和 Summary 类型的指标
- 灵活的配置选项
- 提供 HTTP 和 Gin 框架的中间件

## 安装

使用 Go 模块安装 Everfir Metrics-Go：
```bash
go get github.com/everfir/metrics-go
```

## 使用方法

### 初始化

在使用 Everfir Metrics-Go 之前，您需要先初始化指标系统。您可以选择使用 Collector 模式或 Pushgateway 模式。
```go
import (
    "github.com/everfir/metrics-go"
    "time"
)

// Collector 模式
metrics.Init(metrics.CollectorType, "namespace", "subsystem", 8080, "", "", 0)
// Pushgateway 模式
metrics.Init(metrics.PushgatewayType, "namespace", "subsystem", 0, "http://pushgateway:9091", "job_name", 10time.Second)
```


### 注册指标

在使用指标之前，您需要先注册它们：
```go
import (
    "github.com/everfir/metrics-go"
    "github.com/everfir/metrics-go/structs/metric_info"
)
metrics.Register(metric_info.MetricInfo{
    Type: metric_info.Counter,
    Name: "example_counter",
    Help: "An example counter",
})
metrics.Register(metric_info.MetricInfo{
    Type: metric_info.Gauge,
    Name: "example_gauge",
    Help: "An example gauge",
})
metrics.Register(metric_info.MetricInfo{
    Type: metric_info.Histogram,
    Name: "example_histogram",
    Help: "An example histogram",
    Buckets: []float64{1, 5, 10, 50, 100},
})
```

### 报告指标

使用 `Report` 函数来报告指标值：
```go
import (
    "context"
    "github.com/everfir/metrics-go"
)
metrics.Report(context.Background(), "example_counter", map[string]string{"label": "value"}, 1)
metrics.Report(context.Background(), "example_gauge", map[string]string{"label": "value"}, 50.5)
metrics.Report(context.Background(), "example_histogram", map[string]string{"label": "value"}, 75.0)
```

### 关闭

在应用程序退出时，请确保优雅地关闭指标系统：
```go
import (
    "context"
    "github.com/everfir/metrics-go"
    "time"
)

ctx, cancel := context.WithTimeout(context.Background(), 5time.Second)
defer cancel()
if err := metrics.Close(ctx); err != nil {
    // 处理错误
}
```

## 示例

查看 `example/example.go` 文件以获取完整的使用示例。该示例展示了如何使用命令行参数来选择 Collector 或 Pushgateway 模式，以及如何注册和报告指标。

运行示例：

- Collector 模式：
  ```
  go run example/example.go -mode collector -port 10086
  ```

- Pushgateway 模式：
  ```
  go run example/example.go -mode pushgateway -push-addr http://localhost:9091 -push-interval 5s
  ```

## 项目结构

- `metrics.go`: 主要的 API 实现
- `reporter.go`: 数据上报的具体实现，包含 Collector 和 Pushgateway 两种模式
- `middleware.go`: 中间件，集成在框架中，可以针对 HTTP/RPC/WebSocket 请求中添加指标收集
- `structs/`: 包含各种辅助结构和接口定义
- `example/`: 包含使用示例

## 贡献

欢迎提交问题和拉取请求。对于重大更改，请先开启一个问题讨论您想要更改的内容。

## 许可证

[MIT](https://choosealicense.com/licenses/mit/)