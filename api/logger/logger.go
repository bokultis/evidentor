package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

//Logger replacement for logrus
var Logger *logrus.Logger

//NewLogger replace standard loger with logrus
func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.Formatter = &logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "02-01-2006 15:04:05",
	}

	logger.Out = os.Stdout

	file, err := os.OpenFile("/var/log/logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.Out = file
	} else {
		logger.Info("Failed to log to file, using default stderr")
	}

	Logger = logger
	return Logger
}

//LogLevel custom levels of a log
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
