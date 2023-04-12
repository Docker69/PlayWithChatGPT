// logger.go

package mylogger

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *logrus.Logger

func init() {
	// Create the logger
	Logger = logrus.New()

	Logger.Formatter = new(logrus.JSONFormatter)

	// Set the log level
	Logger.SetLevel(logrus.InfoLevel)

	// Create a new instance of Lumberjack for log rotation
	rotateFileHook := &lumberjack.Logger{
		Filename:   "backend.log",
		MaxSize:    50, // Max size in megabytes before rotation occurs
		MaxBackups: 30, // Max number of old log files to keep
		MaxAge:     2,  // Max number of days to retain log files before deletion
		LocalTime:  true,
	}
	// Set the output of the logger to the file
	Logger.Out = rotateFileHook
}
