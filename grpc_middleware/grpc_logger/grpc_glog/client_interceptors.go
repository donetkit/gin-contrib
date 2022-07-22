package grpc_glog

import (
	"context"
	"encoding/json"
	"time"

	"google.golang.org/grpc"
)

// UnaryClientInterceptor returns a new unary client interceptor that optionally logs the execution of external gRPC calls.
func UnaryClientInterceptor(opts ...Option) grpc.UnaryClientInterceptor {
	o := evaluateClientOpt(opts)
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log := &LogParams{
			TimeStamp: time.Now(),
			IP:        LogRequestIP(ctx),
		}
		LogRequest(log, req, method)
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err != nil {
			LogStatusError(log, err)
		} else {
			LogStatusError(log, err)
		}
		logStr, _ := json.Marshal(log)
		o.logger.Debug(string(logStr))
		return err
	}
}

// StreamClientInterceptor returns a new streaming client interceptor that optionally logs the execution of external gRPC calls.
func StreamClientInterceptor(opts ...Option) grpc.StreamClientInterceptor {
	o := evaluateClientOpt(opts)
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		log := &LogParams{
			TimeStamp: time.Now(),
			IP:        LogRequestIP(ctx),
		}
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			LogStatusError(log, err)
		} else {
			LogStatusError(log, err)
			//LogResponse(log, resp)
		}
		log.Latency = time.Since(log.TimeStamp)
		logStr, _ := json.Marshal(log)
		o.logger.Debug(string(logStr))
		return clientStream, err

	}
}
