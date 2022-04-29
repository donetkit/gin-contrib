package glog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

type ZapLogger struct {
	logger       *zap.SugaredLogger
	config       *Config
	dateFormat   string
	class        string
	logFormatter func(interface{}, bool) string
}

func NewZapLogger(opts ...Option) ILogger {
	cfg := &Config{
		logLevel: DEBUG,
		log2File: false,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	zapLog := &ZapLogger{dateFormat: defaultDateFormat, config: cfg}
	zapLog.SetCustomLogFormat(defaultLogFormatter)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = timeEncodeFunc //指定时间格式
	encoderConfig.TimeKey = "time"
	encoder := zapcore.NewConsoleEncoder(encoderConfig) //获取编码器,NewJSONEncoder()输出json格式，NewConsoleEncoder()输出普通文本格式
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	var writer zapcore.WriteSyncer
	var logLevel zapcore.Level
	switch cfg.logLevel {
	case DEBUG:
		logLevel = zapcore.DebugLevel
	case INFO:
		logLevel = zapcore.InfoLevel
	case WARNING:
		logLevel = zapcore.WarnLevel
	case ERROR:
		logLevel = zapcore.ErrorLevel
	}
	if cfg.log2File {
		fileWriteSyncer := zapcore.AddSync(&Logger{
			Filename:   "./logs/log.log", //日志文件存放目录
			MaxSize:    1,                //文件大小限制,单位MB
			MaxBackups: 15,               //最大保留日志文件数量
			MaxAge:     7,                //日志文件保留天数
			Compress:   false,            //是否压缩处理
		})
		fileCore := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(fileWriteSyncer, zapcore.AddSync(os.Stdout)), logLevel)
		logger := zap.New(fileCore, zap.AddCaller()) //AddCaller()为显示文件名和行号
		zapLog.logger = logger.Sugar()
	} else {
		writer = zapcore.AddSync(os.Stdout)
		logCore := zapcore.NewCore(encoder, writer, logLevel)
		logger := zap.New(logCore)
		zapLog.logger = logger.Sugar()
	}
	return zapLog
}

func timeEncodeFunc(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func (log *ZapLogger) SetClass(className string) {
	log.class = className
}

func (log *ZapLogger) SetCustomLogFormat(logFormatterFunc func(logInfo interface{}, color bool) string) {
	log.logFormatter = logFormatterFunc
}

func (log *ZapLogger) SetDateFormat(format string) {
	log.dateFormat = format
}

func (log *ZapLogger) log(level LogLevel, format string, a ...interface{}) {
	message := format
	message = fmt.Sprintf(format, a...)

	start := time.Now()
	info := LogInfo{
		StartTime: start.Format(log.dateFormat),
		Level:     LevelString[level],
		Class:     log.class,
		Host:      log.config.hostName,
		IP:        log.config.ip,
		Message:   message,
	}
	if log.config.logColor {
		info.Level = LevelColorString[level]
	}
	log.logger.With(log.logFormatter(info, log.config.logColor))
}

func (log *ZapLogger) Debug(format string, a ...interface{}) {
	log.log(DEBUG, format, a...)
}

func (log *ZapLogger) Info(format string, a ...interface{}) {
	log.log(INFO, format, a...)
}

func (log *ZapLogger) Warning(format string, a ...interface{}) {
	log.log(WARNING, format, a...)
}

func (log *ZapLogger) Error(format string, a ...interface{}) {
	log.log(ERROR, format, a...)
}
