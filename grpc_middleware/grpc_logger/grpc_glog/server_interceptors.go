// Copyright (c) Improbable Worlds Ltd, All Rights Reserved

package grpc_glog

import (
	"bytes"
	"context"
	"github.com/donetkit/contrib-gin/grpc_middleware"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"path"
	"strings"
	"time"

	"encoding/json"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor returns a new unary server interceptors that adds logrus.Entry to the context.
func UnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateServerOpt(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := &LogParams{
			TimeStamp: time.Now(),
			IP:        LogRequestIP(ctx),
		}
		LogRequest(log, req, info.FullMethod)
		resp, err := handler(ctx, req)
		if err != nil {
			LogStatusError(log, err)
		} else {
			LogStatusError(log, err)
			LogResponse(log, resp)
		}
		//code := o.codeFunc(err)
		//log.StatusCode = code.String()
		log.Latency = time.Since(log.TimeStamp)
		logStr, _ := json.Marshal(log)
		o.logger.Debug(string(logStr))
		return resp, err
	}
}

// StreamServerInterceptor returns a new streaming server interceptor that adds logrus.Entry to the context.
func StreamServerInterceptor(opts ...Option) grpc.StreamServerInterceptor {
	o := evaluateServerOpt(opts)
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		log := &LogParams{
			TimeStamp: time.Now(),
			IP:        LogRequestIP(ctx),
		}
		wrapped := grpc_middleware.WrapServerStream(stream)
		wrapped.WrappedContext = ctx
		err := handler(srv, wrapped)
		if err != nil {
			LogStatusError(log, err)
		} else {
			LogStatusError(log, err)
			//LogResponse(log, resp)
		}
		log.Latency = time.Since(log.TimeStamp)
		logStr, _ := json.Marshal(log)
		o.logger.Debug(string(logStr))
		return err
	}
}

func LogRequestIP(ctx context.Context) string {
	if p, ok := peer.FromContext(ctx); ok {
		return p.Addr.String()
	}
	return ""
}

func LogRequest(log *LogParams, req interface{}, fullMethodString string) {
	log.Service = path.Dir(fullMethodString)[1:]
	log.Method = path.Base(fullMethodString)
	if b := GetRawJSON(req); b != nil {
		log.RequestData = string(b.Bytes())
	}
}

func LogResponse(log *LogParams, resp interface{}) {
	if b := GetRawJSON(resp); b != nil {
		log.ResponseData = string(b.Bytes())
	}
}

func LogMetadata(log *LogParams, md *metadata.MD) []string {
	var dict []string
	for i := range *md {
		dict = append(dict, i, strings.Join(md.Get(i), ","))
	}
	return dict
}

func LogUserAgent(log *LogParams, md *metadata.MD) {
	if ua := strings.Join(md.Get("user-agent"), ""); ua != "" {
		log.RequestUserAgent = ua
	}
}

func LogStatusError(log *LogParams, err error) {
	statusErr := status.Convert(err)
	//statusErr.Details()
	log.StatusCode = statusErr.Code().String()
	log.ErrorMessage = statusErr.Message()
}

// GetRawJSON converts a Protobuf message to JSON bytes if less than MaxSize.
func GetRawJSON(i interface{}) *bytes.Buffer {
	if pb, ok := i.(proto.Message); ok {
		b := &bytes.Buffer{}
		if err := Marshaller.Marshal(b, pb); err == nil && b.Len() < 2048000 {
			return b
		}
	}
	return nil
}
