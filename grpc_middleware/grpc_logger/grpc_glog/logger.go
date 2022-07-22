package grpc_glog

import (
	"github.com/golang/protobuf/jsonpb"
	"time"
)

var (
	// Marshaller of Protobuf to JSON
	Marshaller = &jsonpb.Marshaler{}
)

// LogParams is the structure any formatter will be handed when time to log comes
type LogParams struct {
	//Service is the HTTP method given to the request.
	Service string
	// Method is the HTTP method given to the request.
	Method string
	// TimeStamp shows the time after the webserve returns a response.
	TimeStamp time.Time
	// StatusCode is HTTP response code.
	StatusCode string
	// Latency is how much time the webserve cost to process a certain request.
	Latency time.Duration
	// IP equals Context's IP method.
	IP string
	// Path is a path the client requests.
	Path string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// isTerm shows whether does gin's output descriptor refers to a terminal.
	isTerm bool
	// BodySize is the size of the Response Body
	BodySize int
	// Keys are the keys set on the request's context.
	Keys map[string]interface{}

	RequestData      string
	RequestUserAgent string
	RequestReferer   string
	RequestProto     string

	RequestId string
	TraceId   string
	SpanId    string

	ResponseData string
}
