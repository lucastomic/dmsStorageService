package middleware

import (
	"net/http"
	"time"

	"github.com/lucastomic/dmsStorageService/internal/logging"
)

// logginMiddleware struct holds a logging.Logger, providing logging capabilities across the application.
// It is designed to log details about HTTP requests processed by the server.
type logginMiddleware struct {
	logger logging.Logger
}

// NewLoggingMiddleware initializes and returns a new instance of logginMiddleware with the provided logging.Logger.
// This function allows for easy creation and integration of logging middleware into the HTTP server's middleware chain.
func NewLoggingMiddleware(logger logging.Logger) Middleware {
	return logginMiddleware{logger}
}

// Execute is the implementation of the Middleware interface for logginMiddleware.
// It wraps an http.HandlerFunc with logging functionality, recording the start time of a request,
// the status code of the response, and the duration of the request processing.
// This method adds a record once the request is processed, recording its information.
func (l logginMiddleware) Execute(
	next http.HandlerFunc,
	errorHandler errorHandler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lwr := newLoggingResponseWriter(w)
		defer func() {
			duration := time.Since(
				start,
			)
			l.logger.Request(r.Context(), r, lwr.statusCode, duration)
		}()
		next(lwr, r)
	}
}
