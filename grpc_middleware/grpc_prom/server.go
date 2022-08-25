// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

// gRPC Prometheus monitoring interceptors for server-side gRPC. github.com/grpc-ecosystem/go-grpc-prometheus

package grpc_prom

import (
	"context"
	"github.com/gin-gonic/gin"
	prom "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

var (
	// DefaultServerMetrics is the default instance of ServerMetrics. It is
	// intended to be used in conjunction the default Prometheus metrics
	// registry.
	DefaultServerMetrics *ServerMetrics

	// UnaryServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Unary RPCs.
	UnaryServerInterceptor func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)

	// StreamServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Streaming RPCs.
	StreamServerInterceptor func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error
)

// RegisterServer takes a gRPC server and pre-initializes all counters to 0. This
// allows for easier monitoring in Prometheus (no missing metrics), and should
// be called *after* all services have been registered with the server. This
// function acts on the DefaultServerMetrics variable.
func RegisterServer(server *grpc.Server, opts ...Option) {
	cfg := &config{
		slowTime:   1,
		namespace:  "service",
		name:       "service",
		duration:   []float64{0.1, 0.3, 1.2, 5},
		handlerUrl: "/metrics",
		endpointLabelMappingFn: func(c *gin.Context) string {
			return c.Request.URL.Path
		},
	}
	for _, opt := range opts {
		opt(cfg)
	}
	DefaultServerMetrics = NewServerMetrics()
	DefaultServerMetrics.config = cfg
	UnaryServerInterceptor = DefaultServerMetrics.UnaryServerInterceptor()
	StreamServerInterceptor = DefaultServerMetrics.StreamServerInterceptor()

	prom.MustRegister(DefaultServerMetrics.serverStartedCounter)
	prom.MustRegister(DefaultServerMetrics.serverHandledCounter)
	prom.MustRegister(DefaultServerMetrics.serverStreamMsgReceived)
	prom.MustRegister(DefaultServerMetrics.serverStreamMsgSent)
	prom.MustRegister(DefaultServerMetrics.serverHandledUptime)

	DefaultServerMetrics.InitializeMetrics(server)
	go DefaultServerMetrics.recordUptime()
}

// EnableHandlingTimeHistogram turns on recording of handling time
// of RPCs. Histogram metrics can be very expensive for Prometheus
// to retain and query. This function acts on the DefaultServerMetrics
// variable and the default Prometheus metrics registry.
func EnableHandlingTimeHistogram(opts ...HistogramOption) {
	DefaultServerMetrics.EnableHandlingTimeHistogram(opts...)
	prom.Register(DefaultServerMetrics.serverHandledHistogram)
}
