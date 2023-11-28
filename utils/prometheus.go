package utils

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func PrometheusBoot(port int) {
	// 在 "/metrics" 路径上注册一个处理器，用于 Prometheus 的数据抓取
	http.Handle("/metrics", promhttp.Handler())

	go func() {
		// 构造监听地址和端口，启动 HTTP 服务
		// 0.0.0.0 表示服务器将接受来自任何 IP 地址的连接，因此可以通过任何可用的 IP 地址和端口来访问该服务器。
		// 监听请求并返回静态内容，使用默认的 HTTP 处理逻辑来处理请求足够了，所以第二个参数nil
		err := http.ListenAndServe("0.0.0.0:"+strconv.Itoa(port), nil)
		// 如果启动失败，记录致命错误并退出
		if err != nil {
			fmt.Println("start fail")
		}
	}()

	// 创建一个新的 CounterVec 指标，用于记录订单请求总数
	orderRequestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_requests_total",           // 指标名称
			Help: "Total number of order requests", // 指标的描述信息
		},
		[]string{"total"}, // 添加一个标签用于标识模块，标签可以用于分类指标
	)
	prometheus.MustRegister(orderRequestsTotal)

	// 创建一个新的 Histogram 指标，用于记录订单请求的处理时间分布
	orderProcessingTime := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "order_processing_time_seconds",                 // 指标名称
			Help: "Histogram of order processing time in seconds", // 指标的描述信息
			// 配置更多的参数，如 Buckets，用于定义处理时间的分档范围
			// 0.1 秒以下的订单数量将归入第一个分档，0.1 到 0.5 秒之间的订单数量将归入第二个分档，以此类推。
			Buckets: []float64{0.1, 0.5, 1, 2, 3},
		},
		[]string{"order_processing"}, // 添加一个标签用于标识模块
	)

	prometheus.MustRegister(orderProcessingTime)

	// 记录日志信息，表明监控服务已启动
	fmt.Println("监控启动，端口为：" + strconv.Itoa(port))
}
