package middleware

import (
	"net/http"
)

// errorHandler is a function type that defines the structure for handling errors within the application.
// It takes an HTTP request, response writer, the error encountered, and the HTTP status code as parameters.
type errorHandler func(*http.Request, http.ResponseWriter, error, int)

// Middleware defines an interface for HTTP middleware components in the application.
// It provides a standard way to process or modify HTTP requests and responses,
// and to handle errors within the HTTP handling pipeline.
type Middleware interface {
	// Execute wraps a given http.HandlerFunc with additional functionality.
	// It takes an http.HandlerFunc and an errorHandler as arguments, and returns a new http.HandlerFunc.
	// The returned http.HandlerFunc is expected to perform any middleware-specific processing,
	// call the next handler, and invoke the errorHandler in case of errors.
	Execute(http.HandlerFunc, errorHandler) http.HandlerFunc
}

// ChainMiddleware takes an http.HandlerFunc, an errorHandler, and a variable number of Middleware interfaces.
// It iteratively wraps the provided http.HandlerFunc with the middleware in reverse order, ensuring the first
// middleware in the slice is the outermost layer and the last is the closest to the original handler.
// This function enables the sequential application of middleware logic to HTTP requests and centralized error handling.
func ChainMiddleware(
	h http.HandlerFunc,
	errorHandler errorHandler,
	middlewares ...Middleware,
) http.HandlerFunc {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i].Execute(h, errorHandler)
	}
	return h
}
