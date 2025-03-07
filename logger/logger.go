package logger

import (
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()
var Entry = logrus.NewEntry(logger)
var WithField = logger.WithField
var Info = logger.Info
var Warn = logger.Warn
var Debug = logger.Debug
var Error = logger.Error

// Alert 需要发送警报的日志
func Alert(args ...interface{}) {
	logger.Error(args...)
}

var Infof = logger.Infof
var Warnf = logger.Warnf
var Debugf = logger.Debugf
var Errorf = logger.Errorf

var Printf = logger.Printf

func GetLogger() *logrus.Logger {
	return logger
}

func init() {
	// if runtime.GOOS != "linux" || os.Getenv("LOGGER_FORMATTER") == "TextFormatter" {
	// 	// 本地调试
	// 	logger.SetFormatter(&logrus.TextFormatter{
	// 		FullTimestamp:   true,
	// 		TimestampFormat: "2006-01-02 15:04:05.000",
	// 		ForceColors:     true,
	// 	})
	// } else {
	// 	logger.SetFormatter(&ecslogrus.Formatter{})
	// 	logger.ReportCaller = true
	// }
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
		ForceColors:     true,
	})

	updateEntry()
}

func updateEntry() {
	WithField = logger.WithField
	Info = logger.Info
	Warn = logger.Warn
	Debug = logger.Debug
	Error = logger.Error

	Infof = logger.Infof
	Warnf = logger.Warnf
	Debugf = logger.Debugf
	Errorf = logger.Errorf
}
