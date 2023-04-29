// logger.go

package utils

import (
	"io"
	"os"

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
		Filename:   "server.log",
		MaxSize:    50, // Max size in megabytes before rotation occurs
		MaxBackups: 30, // Max number of old log files to keep
		MaxAge:     2,  // Max number of days to retain log files before deletion
		LocalTime:  true,
	}
	multi := io.MultiWriter(rotateFileHook, os.Stderr)
	// Set the output of the logger to the file and stdout
	Logger.SetOutput(multi)
}
