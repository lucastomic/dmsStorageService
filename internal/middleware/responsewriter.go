package middleware

import (
	"net/http"
)

// loggingResponseWriter is a wrapper around http.ResponseWriter that captures the status code
// written to the response. This allows middleware or handlers that wrap the response writer
// to log or otherwise act on the HTTP status code after a response has been sent.
type loggingResponseWriter struct {
	http.ResponseWriter     // Embedding the http.ResponseWriter interface.
	statusCode          int // statusCode holds the HTTP status code set by the handler.
}

// newLoggingResponseWriter initializes a new instance of loggingResponseWriter
// wrapping the provided http.ResponseWriter. The default status code is set to http.StatusOK,
// assuming that if no status code is explicitly set, the response will be successful.
func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
}

// WriteHeader captures the HTTP status code set by the handler and delegates the call
// to the underlying http.ResponseWriter's WriteHeader method. This allows the middleware
// to observe and log the status code while still correctly setting the status code
// in the HTTP response.
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code                // Capture the status code for logging or other purposes.
	lrw.ResponseWriter.WriteHeader(code) // Delegate to the original ResponseWriter.
}
