package logging

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/lucastomic/dmsStorageService/internal/contextypes"
	"github.com/sirupsen/logrus"
)

// LogrusLogger is an implementation of the Logger interface using the Logrus logging library.
// It supports logging to multiple destinations, including both terminal and file outputs.
type LogrusLogger struct {
	loggers []*logrus.Logger // A slice of *logrus.Logger instances for logging.
}

// NewLogrusLogger initializes a new LogrusLogger with logging outputs set to both terminal and specified log file.
// It attempts to open or create the log file at the provided logFilePath, logging a fatal error if the file cannot be accessed.
func NewLogrusLogger() Logger {
	ttyLogger := logrus.New()
	return &LogrusLogger{
		loggers: []*logrus.Logger{ttyLogger},
	}
}

// Request logs information about an HTTP request to all configured loggers.
// It formats the log message with details including the timestamp, URI, method, user agent, status code, and duration.
// The RequestID from the context, if present, is included as a field in the log entry.
// For exmaple,
// INFO[0004] 2024-02-03 20:17:12 | /upload | POST | PostmanRuntime/7.36.1 | 201 | 18011  requestID=3
// where 18011 is the duration of the process in microseconds
func (l *LogrusLogger) Request(
	ctx context.Context,
	req *http.Request,
	status int,
	duration time.Duration,
) {
	msg := fmt.Sprintf(
		"%s | %s | %s | %s | %d | %d",
		time.Now().Format("2006-01-02 15:04:05"),
		req.RequestURI,
		req.Method,
		req.UserAgent(),
		status,
		duration.Microseconds(),
	)

	for _, logger := range l.loggers {
		logger.WithField("requestID", ctx.Value(contextypes.CTXRequestIDKey{})).Info(msg)
	}
}

// Info logs an informational message to all configured loggers.
// The message is formatted according to the provided format string and arguments.
// The RequestID from the context, if present, is included as a field in the log entry.
func (l *LogrusLogger) Info(ctx context.Context, format string, a ...any) {
	for _, logger := range l.loggers {
		logger.WithField("requestID", ctx.Value(contextypes.CTXRequestIDKey{})).
			Info(fmt.Sprintf(format, a...))
	}
}

// Error logs an error message to all configured loggers.
// The message is formatted according to the provided format string and arguments.
// The RequestID from the context, if present, is included as a field in the log entry.
func (l *LogrusLogger) Error(ctx context.Context, format string, a ...any) {
	for _, logger := range l.loggers {
		logger.WithField("requestID", ctx.Value(contextypes.CTXRequestIDKey{})).
			Error(fmt.Sprintf(format, a...))
	}
}
