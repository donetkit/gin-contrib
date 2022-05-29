package prom

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"regexp"
	"time"
)

var (
	labels = []string{"status", "endpoint", "method"}

	labelsServeName = []string{"name"}

	reqUVTotal *prometheus.CounterVec

	slowReqTotal *prometheus.CounterVec

	uptime *prometheus.CounterVec

	reqCount *prometheus.CounterVec

	reqDuration *prometheus.HistogramVec

	reqSizeBytes *prometheus.SummaryVec

	respSizeBytes *prometheus.SummaryVec
)

// init registers the prometheus metrics
func (c *config) registerPrometheusOpts() {

	reqUVTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: c.namespace,
			Name:      "request_uv_total",
			Help:      "all the server received ip num.",
		}, nil,
	)

	slowReqTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: c.namespace,
			Name:      "slow_request_total",
			Help:      fmt.Sprintf("the server handled slow requests counter, t=%d.", int(c.slowTime)),
		}, labels,
	)

	uptime = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: c.namespace,
			Name:      "uptime",
			Help:      "HTTP service uptime.",
		}, labelsServeName,
	)

	reqCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: c.namespace,
			Name:      "http_request_count_total",
			Help:      "Total number of HTTP requests made.",
		}, labels,
	)

	reqDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: c.namespace,
			Name:      "http_request_duration_seconds",
			Help:      "HTTP request latencies in seconds.",
			Buckets:   c.duration,
		}, labels,
	)

	reqSizeBytes = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: c.namespace,
			Name:      "http_request_size_bytes",
			Help:      "HTTP request sizes in bytes.",
		}, labels,
	)

	respSizeBytes = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: c.namespace,
			Name:      "http_response_size_bytes",
			Help:      "HTTP response sizes in bytes.",
		}, labels,
	)
	prometheus.MustRegister(reqUVTotal, slowReqTotal, uptime, reqCount, reqDuration, reqSizeBytes, respSizeBytes)
	go c.recordUptime()
}

// recordUptime increases service uptime per second.
func (c *config) recordUptime() {
	for range time.Tick(time.Second) {
		uptime.WithLabelValues(c.name).Inc()
	}
}

// calcRequestSize returns the size of request object.
func calcRequestSize(r *http.Request) float64 {
	size := 0
	if r.URL != nil {
		size = len(r.URL.String())
	}

	size += len(r.Method)
	size += len(r.Proto)

	for name, values := range r.Header {
		size += len(name)
		for _, value := range values {
			size += len(value)
		}
	}
	size += len(r.Host)

	// r.Form and r.MultipartForm are assumed to be included in r.URL.
	if r.ContentLength != -1 {
		size += int(r.ContentLength)
	}
	return float64(size)
}

type RequestLabelMappingFn func(c *gin.Context) string

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

// New returns a gin.HandlerFunc for exporting some Web metrics
func New(opts ...Option) gin.HandlerFunc {
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
	cfg.registerPrometheusOpts()
	bloomFilter := NewBloomFilter()
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := fmt.Sprintf("%d", c.Writer.Status())
		endpoint := cfg.endpointLabelMappingFn(c)
		method := c.Request.Method

		lvs := []string{status, endpoint, method}

		isOk := cfg.checkLabel(status, cfg.excludeRegexStatus) && cfg.checkLabel(endpoint, cfg.excludeRegexEndpoint) && cfg.checkLabel(method, cfg.excludeRegexMethod)

		if !isOk {
			return
		}
		// no response content will return -1
		respSize := c.Writer.Size()
		if respSize < 0 {
			respSize = 0
		}

		// set uv
		if clientIP := c.ClientIP(); !bloomFilter.Contains(clientIP) {
			bloomFilter.Add(clientIP)
			reqUVTotal.WithLabelValues().Inc()
		}

		second := time.Since(start).Seconds()
		
		// set slow request
		if second > cfg.slowTime {
			slowReqTotal.WithLabelValues(lvs...).Inc()
		}
		reqCount.WithLabelValues(lvs...).Inc()
		reqDuration.WithLabelValues(lvs...).Observe(second)
		reqSizeBytes.WithLabelValues(lvs...).Observe(calcRequestSize(c.Request))
		respSizeBytes.WithLabelValues(lvs...).Observe(float64(respSize))
	}
}

// promHandler wrappers the standard http.Handler to gin.HandlerFunc
func promHandler(handler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
