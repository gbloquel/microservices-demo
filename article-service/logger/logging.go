package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger is the global instance of the logrus.Logger used for logging in the application.
var Logger *logrus.Logger = logrus.New()

// SetupLogging initializes the global Logger with the specified log level and JSON formatter.
// It reads the log level configuration from the environment variable LOG_LEVEL.
// If LOG_LEVEL is not set or is invalid, it defaults to "info" log level.
// The timestamp format for log entries is set to "2006-01-02 15:04:05Z".
func SetupLogging() {
	// Read log level configuration from environment variable LOG_LEVEL
	logLevel := os.Getenv("LOG_LEVEL")

	// Set default log level to info if LOG_LEVEL is not set or is invalid
	if logLevel == "" {
		logLevel = "info"
	}

	// Parse the log level and set it for the global Logger
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	// Set the formatter for the global Logger to JSON format with a specific timestamp format
	Logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05Z"})
}
