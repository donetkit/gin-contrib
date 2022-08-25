// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

// gRPC Prometheus monitoring interceptors for client-side gRPC.

package grpc_prom

import (
	"context"
	"github.com/gin-gonic/gin"
	prom "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

var (
	// DefaultClientMetrics is the default instance of ClientMetrics. It is
	// intended to be used in conjunction the default Prometheus metrics
	// registry.
	DefaultClientMetrics *ClientMetrics

	// UnaryClientInterceptor is a gRPC client-side interceptor that provides Prometheus monitoring for Unary RPCs.
	UnaryClientInterceptor func(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error

	// StreamClientInterceptor is a gRPC client-side interceptor that provides Prometheus monitoring for Streaming RPCs.
	StreamClientInterceptor func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error)
)

// RegisterClient takes a gRPC server and pre-initializes all counters to 0. This
// allows for easier monitoring in Prometheus (no missing metrics), and should
// be called *after* all services have been registered with the server. This
// function acts on the DefaultClientMetrics variable.
func RegisterClient(server *grpc.Server, opts ...Option) {
	cfg := &config{
		slowTime:   1,
		namespace:  "client",
		name:       "client",
		duration:   []float64{0.1, 0.3, 1.2, 5},
		handlerUrl: "/metrics",
		endpointLabelMappingFn: func(c *gin.Context) string {
			return c.Request.URL.Path
		},
	}
	for _, opt := range opts {
		opt(cfg)
	}

	DefaultClientMetrics = NewClientMetrics()
	DefaultClientMetrics.config = cfg
	UnaryClientInterceptor = DefaultClientMetrics.UnaryClientInterceptor()
	StreamClientInterceptor = DefaultClientMetrics.StreamClientInterceptor()

	prom.MustRegister(DefaultClientMetrics.clientStartedCounter)
	prom.MustRegister(DefaultClientMetrics.clientHandledCounter)
	prom.MustRegister(DefaultClientMetrics.clientStreamMsgReceived)
	prom.MustRegister(DefaultClientMetrics.clientStreamMsgSent)

	DefaultServerMetrics.InitializeMetrics(server)
	go DefaultServerMetrics.recordUptime()
}

// EnableClientHandlingTimeHistogram turns on recording of handling time of
// RPCs. Histogram metrics can be very expensive for Prometheus to retain and
// query. This function acts on the DefaultClientMetrics variable and the
// default Prometheus metrics registry.
func EnableClientHandlingTimeHistogram(opts ...HistogramOption) {
	DefaultClientMetrics.EnableClientHandlingTimeHistogram(opts...)
	prom.Register(DefaultClientMetrics.clientHandledHistogram)
}

// EnableClientStreamReceiveTimeHistogram turns on recording of
// single message receive time of streaming RPCs.
// This function acts on the DefaultClientMetrics variable and the
// default Prometheus metrics registry.
func EnableClientStreamReceiveTimeHistogram(opts ...HistogramOption) {
	DefaultClientMetrics.EnableClientStreamReceiveTimeHistogram(opts...)
	prom.Register(DefaultClientMetrics.clientStreamRecvHistogram)
}

// EnableClientStreamSendTimeHistogram turns on recording of
// single message send time of streaming RPCs.
// This function acts on the DefaultClientMetrics variable and the
// default Prometheus metrics registry.
func EnableClientStreamSendTimeHistogram(opts ...HistogramOption) {
	DefaultClientMetrics.EnableClientStreamSendTimeHistogram(opts...)
	prom.Register(DefaultClientMetrics.clientStreamSendHistogram)
}
