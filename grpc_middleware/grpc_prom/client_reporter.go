// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package grpc_prom

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
)

type clientReporter struct {
	metrics     *ClientMetrics
	rpcType     grpcType
	serviceName string
	methodName  string
	startTime   time.Time
	lvs         []string
}

func newClientReporter(m *ClientMetrics, rpcType grpcType, fullMethod string) *clientReporter {
	r := &clientReporter{
		metrics: m,
		rpcType: rpcType,
	}

	r.serviceName, r.methodName = splitMethodName(fullMethod)
	isOk := m.checkLabel(string(r.rpcType), m.config.excludeRegexRpcType) && m.checkLabel(r.serviceName, m.config.excludeRegexServiceName) && m.checkLabel(r.methodName, m.config.excludeRegexMethodName)
	if !isOk {
		return r
	}
	r.lvs = []string{string(r.rpcType), r.serviceName, r.methodName}
	if r.metrics.clientHandledHistogramEnabled {
		r.startTime = time.Now()
	}
	r.serviceName, r.methodName = splitMethodName(fullMethod)
	r.metrics.clientStartedCounter.WithLabelValues(r.lvs...).Inc()
	return r
}

// timer is a helper interface to time functions.
type timer interface {
	ObserveDuration() time.Duration
}

type noOpTimer struct {
}

func (noOpTimer) ObserveDuration() time.Duration {
	return 0
}

var emptyTimer = noOpTimer{}

func (r *clientReporter) ReceiveMessageTimer() timer {
	isOk := r.metrics.checkLabel(string(r.rpcType), r.metrics.config.excludeRegexRpcType) && r.metrics.checkLabel(r.serviceName, r.metrics.config.excludeRegexServiceName) && r.metrics.checkLabel(r.methodName, r.metrics.config.excludeRegexMethodName)
	if !isOk {
		return emptyTimer
	}
	if r.metrics.clientStreamRecvHistogramEnabled {
		hist := r.metrics.clientStreamRecvHistogram.WithLabelValues(r.lvs...)
		return prometheus.NewTimer(hist)
	}

	return emptyTimer
}

func (r *clientReporter) ReceivedMessage() {
	isOk := r.metrics.checkLabel(string(r.rpcType), r.metrics.config.excludeRegexRpcType) && r.metrics.checkLabel(r.serviceName, r.metrics.config.excludeRegexServiceName) && r.metrics.checkLabel(r.methodName, r.metrics.config.excludeRegexMethodName)
	if !isOk {
		return
	}
	r.metrics.clientStreamMsgReceived.WithLabelValues(r.lvs...).Inc()
}

func (r *clientReporter) SendMessageTimer() timer {
	isOk := r.metrics.checkLabel(string(r.rpcType), r.metrics.config.excludeRegexRpcType) && r.metrics.checkLabel(r.serviceName, r.metrics.config.excludeRegexServiceName) && r.metrics.checkLabel(r.methodName, r.metrics.config.excludeRegexMethodName)
	if !isOk {
		return emptyTimer
	}
	if r.metrics.clientStreamSendHistogramEnabled {
		hist := r.metrics.clientStreamSendHistogram.WithLabelValues(r.lvs...)
		return prometheus.NewTimer(hist)
	}

	return emptyTimer
}

func (r *clientReporter) SentMessage() {
	isOk := r.metrics.checkLabel(string(r.rpcType), r.metrics.config.excludeRegexRpcType) && r.metrics.checkLabel(r.serviceName, r.metrics.config.excludeRegexServiceName) && r.metrics.checkLabel(r.methodName, r.metrics.config.excludeRegexMethodName)
	if !isOk {
		return
	}
	r.metrics.clientStreamMsgSent.WithLabelValues(r.lvs...).Inc()
}

func (r *clientReporter) Handled(code codes.Code) {
	isOk := r.metrics.checkLabel(code.String(), r.metrics.config.excludeRegexCode) && r.metrics.checkLabel(string(r.rpcType), r.metrics.config.excludeRegexRpcType) && r.metrics.checkLabel(r.serviceName, r.metrics.config.excludeRegexServiceName) && r.metrics.checkLabel(r.methodName, r.metrics.config.excludeRegexMethodName)
	if !isOk {
		return
	}
	r.metrics.clientHandledCounter.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName, code.String()).Inc()
	if r.metrics.clientHandledHistogramEnabled {
		r.metrics.clientHandledHistogram.WithLabelValues(r.lvs...).Observe(time.Since(r.startTime).Seconds())
	}
}
