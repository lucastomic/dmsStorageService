package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRequestIDMiddlewareWithHeader tests the requestIDMiddleware ensuring it passes the request through
// when the X-Request-ID header is present.
func TestRequestIDMiddlewareWithHeader(t *testing.T) {
	middleware := NewRequestIDMiddleware()
	nextCalled := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	errorHandler := func(r *http.Request, w http.ResponseWriter, err error, statusCode int) {
		t.Errorf("errorHandler should not be called when X-Request-ID is present")
	}

	handlerToTest := middleware.Execute(next, errorHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", "test-id")
	w := httptest.NewRecorder()

	handlerToTest.ServeHTTP(w, req)

	if !nextCalled {
		t.Errorf("Next handler was not called")
	}
}

// TestRequestIDMiddlewareWithoutHeader tests the requestIDMiddleware ensuring it calls the errorHandler
// when the X-Request-ID header is missing.
func TestRequestIDMiddlewareWithoutHeader(t *testing.T) {
	middleware := NewRequestIDMiddleware()
	nextCalled := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	errorHandlerCalled := false

	errorHandler := func(r *http.Request, w http.ResponseWriter, err error, statusCode int) {
		if statusCode != http.StatusBadRequest {
			t.Errorf("Expected status code 400, got %v", statusCode)
		}
		if err.Error() != "X-Request-ID can't be null" {
			t.Errorf("Expected error message 'X-Request-ID can't be null', got '%v'", err)
		}
		errorHandlerCalled = true
	}

	handlerToTest := middleware.Execute(next, errorHandler)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handlerToTest.ServeHTTP(w, req)

	if nextCalled {
		t.Errorf("Next handler should not be called when X-Request-ID is missing")
	}

	if !errorHandlerCalled {
		t.Errorf("Error handler was not called when X-Request-ID is missing")
	}
}
