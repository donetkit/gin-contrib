// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grpc_trace // import "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

import (
	"context"
	"google.golang.org/grpc/metadata"
	"regexp"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const (
	// GRPCStatusCodeKey is convention for numeric status code of a gRPC request.
	GRPCStatusCodeKey = attribute.Key("rpc.grpc.status_code")
)

type metadataSupplier struct {
	metadata *metadata.MD
}

// assert that metadataSupplier implements the TextMapCarrier interface
var _ propagation.TextMapCarrier = &metadataSupplier{}

func (s *metadataSupplier) Get(key string) string {
	values := s.metadata.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (s *metadataSupplier) Set(key string, value string) {
	s.metadata.Set(key, value)
}

func (s *metadataSupplier) Keys() []string {
	out := make([]string, 0, len(*s.metadata))
	for key := range *s.metadata {
		out = append(out, key)
	}
	return out
}

// Inject injects correlation context and span context into the gRPC
// metadata object. This function is meant to be used on outgoing
// requests.
func Inject(ctx context.Context, metadata *metadata.MD, opts ...Option) {
	c := newConfig(opts)
	c.tracerServer.Propagators.Inject(ctx, &metadataSupplier{
		metadata: metadata,
	})
}

// Extract returns the correlation context and span context that
// another service encoded in the gRPC metadata object with Inject.
// This function is meant to be used on incoming requests.
func Extract(ctx context.Context, metadata *metadata.MD, opts ...Option) (baggage.Baggage, trace.SpanContext) {
	c := newConfig(opts)
	ctx = c.tracerServer.Propagators.Extract(ctx, &metadataSupplier{
		metadata: metadata,
	})

	return baggage.FromContext(ctx), trace.SpanContextFromContext(ctx)
}

// checkLabel returns the match result of labels.
// Return true if regex-pattern compiles failed.
func (c *config) checkLabel(label string, patterns []string) bool {
	if len(patterns) <= 0 {
		return true
	}
	for _, pattern := range patterns {
		if pattern == "" {
			return true
		}
		matched, err := regexp.MatchString(pattern, label)
		if err != nil {
			return true
		}
		if matched {
			return false
		}
	}
	return true
}
