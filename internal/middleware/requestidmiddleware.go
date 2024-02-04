package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/lucastomic/dmsStorageService/internal/contextypes"
)

// requestIDMiddleware is a middleware that ensures each HTTP request contains an X-Request-ID header.
// If a request does not include this header, the middleware invokes the provided errorHandler
// with a BadRequest status, indicating that the X-Request-ID header is required.
type requestIDMiddleware struct{}

// NewRequestIDMiddleware creates and returns a new instance of requestIDMiddleware.
// This middleware can be used to enforce the presence of an X-Request-ID header in incoming HTTP requests.
func NewRequestIDMiddleware() Middleware {
	return requestIDMiddleware{}
}

// Execute wraps the next http.HandlerFunc in the middleware chain,
// checking for the presence of an X-Request-ID header in the request.
// If the header is missing, it calls the errorHandler with a 400 status code and an appropriate error message.
// Otherwise, it adds the request ID to the request's context and proceeds with the next handler.
func (requestIDMiddleware) Execute(
	next http.HandlerFunc,
	errorHandler errorHandler,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			errorHandler(r, w, errors.New("X-Request-ID can't be null"), http.StatusBadRequest)
			return
		}
		ctx := context.WithValue(r.Context(), contextypes.CTXRequestIDKey{}, requestID)
		*r = *r.WithContext(ctx)
		next(w, r)
	}
}
