package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lucastomic/dmsStorageService/internal/contextypes"
)

// pathVarsMiddleware is a middleware component that extracts variables from the route
// using the gorilla/mux package and adds them to the context of each HTTP request.
// This allows subsequent handlers in the chain to access route variables directly from the request context,
// promoting a decoupled architecture where handlers are not directly dependent on the routing mechanism.
type pathVarsMiddleware struct{}

// NewPathVarsMiddleware creates and returns a new instance of varsMiddleware.
// This function serves as a constructor for varsMiddleware, encapsulating the creation logic
// and ensuring that new instances are properly initialized.
func NewPathVarsMiddleware() Middleware {
	return pathVarsMiddleware{}
}

// Execute is the implementation of the Middleware interface for varsMiddleware.
// It wraps an http.HandlerFunc with additional logic to extract route variables
// and insert them into the request context.
func (pathVarsMiddleware) Execute(
	next http.HandlerFunc,
	errorHandler errorHandler,
) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ctx := r.Context()
		for key, value := range vars {
			ctx = context.WithValue(ctx, contextypes.ContextPathVarKey(key), value)
		}
		*r = *r.WithContext(ctx)
		next(w, r)
	})
}
