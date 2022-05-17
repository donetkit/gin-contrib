package logger

import (
	"github.com/donetkit/contrib-log/glog"
	"github.com/donetkit/gin-contrib/utils/buffer"
	"github.com/gin-gonic/gin"
	"strings"
)

func SetGinDefaultWriter(logger glog.ILogger) {
	gin.DefaultWriter = &writeLogger{pool: buffer.Pool{}, logger: logger}
}

type writeLogger struct {
	logger glog.ILogger
	pool   buffer.Pool
}

// Write implements io.Writer.
func (l *writeLogger) Write(p []byte) (n int, err error) {
	buf := l.pool.Get()
	defer l.pool.Put(buf)
	n, err = buf.Write(p)
	if l.logger != nil {
		msg := buf.String()
		if strings.HasSuffix(msg, "\n") {
			msg = msg[:len(msg)-2]
		}
		l.logger.Info(msg)
	}
	return n, err
}
