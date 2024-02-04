package logging

import (
	"context"
	"net/http"
	"time"
)

// Logger defines the interface for logging throughout the application.
// It provides structured logging capabilities for various levels of application events,
// including requests, informational messages, and errors.
type Logger interface {
	// Request logs an HTTP request event, including details such as the request method, URL,
	// the response status code, and the duration of the request handling.
	Request(ctx context.Context, req *http.Request, statusCode int, duration time.Duration)

	// Info logs an informational message. This method is intended for logging general,
	// non-critical information about application operation. The message format and arguments
	// are similar to fmt.Printf, allowing for flexible message construction.
	Info(ctx context.Context, format string, a ...any)

	// Error logs an error event. This method is used for logging errors and exceptions
	// that occur during application execution. Like Info, it supports formatted messages,
	// making it suitable for reporting issues with detailed context.
	Error(ctx context.Context, format string, a ...any)
}
