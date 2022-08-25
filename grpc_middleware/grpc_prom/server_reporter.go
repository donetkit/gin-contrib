// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package grpc_prom

import (
	"time"

	"google.golang.org/grpc/codes"
)

type serverReporter struct {
	metrics     *ServerMetrics
	rpcType     grpcType
	serviceName string
	methodName  string
	startTime   time.Time
	lvs         []string
}

func newServerReporter(m *ServerMetrics, rpcType grpcType, fullMethod string) *serverReporter {
	r := &serverReporter{
		metrics: m,
		rpcType: rpcType,
	}

	r.serviceName, r.methodName = splitMethodName(fullMethod)
	isOk := m.checkLabel(string(r.rpcType), m.config.excludeRegexRpcType) && m.checkLabel(r.serviceName, m.config.excludeRegexServiceName) && m.checkLabel(r.methodName, m.config.excludeRegexMethodName)
	if !isOk {
		return r
	}

	r.lvs = []string{string(r.rpcType), r.serviceName, r.methodName}
	if r.metrics.serverHandledHistogramEnabled {
		r.startTime = time.Now()
	}
	r.metrics.serverStartedCounter.WithLabelValues(r.lvs...).Inc()
	return r
}

func (r *serverReporter) ReceivedMessage() {
	isOk := r.metrics.checkLabel(string(r.rpcType), r.metrics.config.excludeRegexRpcType) && r.metrics.checkLabel(r.serviceName, r.metrics.config.excludeRegexServiceName) && r.metrics.checkLabel(r.methodName, r.metrics.config.excludeRegexMethodName)
	if !isOk {
		return
	}

	r.metrics.serverStreamMsgReceived.WithLabelValues(r.lvs...).Inc()
}

func (r *serverReporter) SentMessage() {
	isOk := r.metrics.checkLabel(string(r.rpcType), r.metrics.config.excludeRegexRpcType) && r.metrics.checkLabel(r.serviceName, r.metrics.config.excludeRegexServiceName) && r.metrics.checkLabel(r.methodName, r.metrics.config.excludeRegexMethodName)
	if !isOk {
		return
	}
	r.metrics.serverStreamMsgSent.WithLabelValues().Inc()
}

func (r *serverReporter) Handled(code codes.Code) {
	isOk := r.metrics.checkLabel(code.String(), r.metrics.config.excludeRegexCode) && r.metrics.checkLabel(string(r.rpcType), r.metrics.config.excludeRegexRpcType) && r.metrics.checkLabel(r.serviceName, r.metrics.config.excludeRegexServiceName) && r.metrics.checkLabel(r.methodName, r.metrics.config.excludeRegexMethodName)
	if !isOk {
		return
	}
	r.metrics.serverHandledCounter.WithLabelValues(string(r.rpcType), r.serviceName, r.methodName, code.String()).Inc()
	if r.metrics.serverHandledHistogramEnabled {
		r.metrics.serverHandledHistogram.WithLabelValues(r.lvs...).Observe(time.Since(r.startTime).Seconds())
	}
}
