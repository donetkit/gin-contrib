package glog

import (
	"fmt"
	"github.com/donetkit/gin-contrib/glog/console_colors"
)

type ILogger interface {
	SetClass(className string)
	Debug(format string, a ...interface{})
	Info(format string, a ...interface{})
	Warning(format string, a ...interface{})
	Error(format string, a ...interface{})

	SetCustomLogFormat(logFormatterFunc func(logInfo interface{}, color bool) string)
	SetDateFormat(format string)
}

type Config struct {
	log2File bool
	hostName string
	ip       string
	logLevel LogLevel
	logColor bool
}

type LogInfo struct {
	StartTime string
	Level     string
	Class     string
	Host      string
	IP        string
	Message   string
	Extend    map[string]interface{}
}

var defaultDateFormat = "2006/01/02 15:04:05.00.000"

type LogLevel int

// MessageLevel
const (
	DEBUG   = LogLevel(0) // DEBUG = 0
	INFO    = LogLevel(1) // INFO = 1
	WARNING = LogLevel(2) // WARNING = 2
	ERROR   = LogLevel(3) // ERROR = 3
)

var LevelString = map[LogLevel]string{
	DEBUG:   "DEBUG",
	INFO:    "INFO ",
	WARNING: "WARN ",
	ERROR:   "ERROR",
}

var LevelColorString = map[LogLevel]string{
	DEBUG:   console_colors.Green("DEBUG"),
	INFO:    "INFO ",
	WARNING: console_colors.Yellow("WARN "),
	ERROR:   console_colors.Red("ERROR"),
}

func defaultLogFormatter(log interface{}, color bool) string {
	logInfo := log.(LogInfo)
	if color {
		if logInfo.Class != "" {
			logInfo.Class = fmt.Sprintf(" - [%s]", console_colors.Green(logInfo.Class))
		}
		if logInfo.Host != "" && logInfo.IP != "" {
			logInfo.Host = fmt.Sprintf(" - [%s(%s)]", console_colors.Green(logInfo.Host), console_colors.Yellow(logInfo.IP))
		}
		return fmt.Sprintf("%s%s%s - [%s] %s", logInfo.StartTime, logInfo.Host, logInfo.Class, logInfo.Level, logInfo.Message)
	}
	if logInfo.Class != "" {
		logInfo.Class = fmt.Sprintf(" - [%s]", logInfo.Class)
	}
	if logInfo.Host != "" && logInfo.IP != "" {
		logInfo.Host = fmt.Sprintf(" - [%s(%s)]", logInfo.Host, logInfo.IP)
	}
	return fmt.Sprintf("%s%s%s - [%s] %s", logInfo.StartTime, logInfo.Host, logInfo.Class, logInfo.Level, logInfo.Message)
}
