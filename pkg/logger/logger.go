package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logrus.New()

	// Set output to stdout
	log.SetOutput(os.Stdout)

	// Set format based on LOG_FORMAT env
	format := strings.ToLower(os.Getenv("LOG_FORMAT"))
	if format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}

	// Set level based on LOG_LEVEL env
	levelStr := strings.ToLower(os.Getenv("LOG_LEVEL"))
	level, err := logrus.ParseLevel(levelStr)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)
}

// GetLogger returns the global logger instance
func GetLogger() *logrus.Logger {
	return log
}
