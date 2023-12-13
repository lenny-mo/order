package utils

import (
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	
)

// 创建链路追踪实例
func NewTracer(servicename string, addr string) (opentracing.Tracer, io.Closer, error) {
	cfg := config.Configuration{
		ServiceName: servicename,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst, // 固定采样
			Param: 1,                       // 1=全采样, 采样率为 100%，即所有请求都会被采样
		},
		Reporter: &config.ReporterConfig{
			BufferFlushInterval: 1 * time.Second, // 发送间隔
			LocalAgentHostPort:  addr,            // agent 地址
			LogSpans:            true,            // 是否打印日志
		},
	}

	return cfg.NewTracer()
}
