package grpc_glog

import (
	"bytes"
	"context"
	"fmt"
	"github.com/donetkit/contrib-gin/grpc_middleware/grpc_logger"
	"github.com/donetkit/contrib-log/glog"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	// JsonPbMarshaller is the marshaller used for serializing protobuf messages.
	// If needed, this variable can be reassigned with a different marshaller with the same Marshal() signature.
	JsonPbMarshaller grpc_logger.JsonPbMarshaler = &jsonpb.Marshaler{}
)

// PayloadUnaryServerInterceptor returns a new unary server interceptors that logs the payloads of requests.
//
// This *only* works when placed *after* the `grpc_logrus.UnaryServerInterceptor`. However, the logging can be done to a
// separate instance of the logger.
func PayloadUnaryServerInterceptor(opts ...Option) grpc.UnaryServerInterceptor {
	o := evaluateServerOpt(opts)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logProtoMessageAsJsonRequest(o.logger, req)
		resp, err := handler(ctx, req)
		if err == nil {
			logProtoMessageAsJsonResponse(o.logger, resp)
		}
		return resp, err
	}
}

// PayloadStreamServerInterceptor returns a new server server interceptors that logs the payloads of requests.
//
// This *only* works when placed *after* the `grpc_logrus.StreamServerInterceptor`. However, the logging can be done to a
// separate instance of the logger.
func PayloadStreamServerInterceptor(opts ...Option) grpc.StreamServerInterceptor {
	o := evaluateServerOpt(opts)
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// Use the provided logrus.Entry for logging but use the fields from context.
		//logEntry := logger.WithFields(ctxlogrus.Extract(stream.Context()).Data)
		newStream := &loggingServerStream{ServerStream: stream, logger: o.logger}
		return handler(srv, newStream)
	}
}

// PayloadUnaryClientInterceptor returns a new unary client interceptor that logs the payloads of requests and responses.
func PayloadUnaryClientInterceptor(opts ...Option) grpc.UnaryClientInterceptor {
	o := evaluateServerOpt(opts)
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		logProtoMessageAsJsonRequest(o.logger, req)
		err := invoker(ctx, method, req, reply, cc, opts...)
		if err == nil {
			logProtoMessageAsJsonResponse(o.logger, req)
		}
		return err
	}
}

// PayloadStreamClientInterceptor returns a new streaming client interceptor that logs the payloads of requests and responses.
func PayloadStreamClientInterceptor(opts ...Option) grpc.StreamClientInterceptor {
	o := evaluateServerOpt(opts)
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		newStream := &loggingClientStream{ClientStream: clientStream, logger: o.logger}
		return newStream, err
	}
}

type loggingClientStream struct {
	grpc.ClientStream
	logger glog.ILoggerEntry
}

func (l *loggingClientStream) SendMsg(m interface{}) error {
	err := l.ClientStream.SendMsg(m)
	if err == nil {
		logProtoMessageAsJsonRequest(l.logger, m)
	}

	return err
}

func (l *loggingClientStream) RecvMsg(m interface{}) error {
	err := l.ClientStream.RecvMsg(m)
	if err == nil {
		logProtoMessageAsJsonResponse(l.logger, m)
	}
	return err
}

type loggingServerStream struct {
	grpc.ServerStream
	logger glog.ILoggerEntry
}

func (l *loggingServerStream) SendMsg(m interface{}) error {
	err := l.ServerStream.SendMsg(m)
	if err == nil {
		logProtoMessageAsJsonResponse(l.logger, m)
	}
	return err
}

func (l *loggingServerStream) RecvMsg(m interface{}) error {
	err := l.ServerStream.RecvMsg(m)
	if err == nil {
		logProtoMessageAsJsonRequest(l.logger, m)
	}
	return err
}

func logProtoMessageAsJsonRequest(logger glog.ILoggerEntry, pbMsg interface{}) {
	if p, ok := pbMsg.(proto.Message); ok {
		if log, okLog := logger.(*logrus.Entry); okLog {
			log.WithField("Request", "Request").Debug(&jsonpbMarshalleble{p})
		}
	}
}

func logProtoMessageAsJsonResponse(logger glog.ILoggerEntry, pbMsg interface{}) {
	if p, ok := pbMsg.(proto.Message); ok {
		if log, okLog := logger.(*logrus.Entry); okLog {
			log.WithField("Response", "Response").Debug(&jsonpbMarshalleble{p})
		}

	}
}

type jsonpbMarshalleble struct {
	proto.Message
}

func (j *jsonpbMarshalleble) MarshalJSON() ([]byte, error) {
	b := &bytes.Buffer{}
	if err := JsonPbMarshaller.Marshal(b, j.Message); err != nil {
		return nil, fmt.Errorf("jsonpb serializer failed: %v", err)
	}
	return b.Bytes(), nil
}
