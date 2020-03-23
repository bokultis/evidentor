package logger

import (
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.Formatter = &logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "02-01-2006 15:04:05",
	}

	Logger = logger
	return Logger
}

func LogLevel(lvl string) logrus.Level {
	switch lvl {
	case "info":
		return logrus.InfoLevel
	case "error":
		return logrus.ErrorLevel
	default:
		panic("Not supported")
	}
}
