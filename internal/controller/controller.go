package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/lucastomic/dmsStorageService/internal/controller/apitypes"
	"github.com/lucastomic/dmsStorageService/internal/errs"
)

// Controller defines the interface for an HTTP controller.
// It outlines the necessary functionality for routing HTTP requests to their respective handlers.
type Controller interface {
	// Router returns a Router which maps HTTP paths to API functions.
	// This API functions act like the handler for every route.
	Router() apitypes.Router
}

// CommonController provides shared functionalities for handling HTTP requests across different controllers.
type CommonController struct{}

// ParseError takes an error and maps it to an HTTP response.
// It uses the mapDomainErrorToHTTP function to convert domain-specific errors to HTTP errors,
// encapsulating the process of translating backend errors into user-friendly HTTP responses.
// It's important to take into account that if an error is not already defined into the domain-specific errors,
// it will map the error as a internalServerError
func (c *CommonController) ParseError(
	ctx context.Context,
	r *http.Request,
	w http.ResponseWriter,
	err error,
) apitypes.Response {
	httpErr := mapDomainErrorToHTTP(err)
	return apitypes.Response{
		Status:  httpErr.Code,
		Content: map[string]any{"error": httpErr.Error()},
		Headers: map[string]string{"Content-Type": "application/json"},
	}
}

// mapDomainErrorToHTTP converts domain-specific errors to errs.HTTPError instances.
// It checks for specific known errors and maps them to appropriate HTTP status codes and messages.
// For unrecognized errors, it defaults to returning an "internal error" with a 500 status code.
func mapDomainErrorToHTTP(err error) errs.HTTPError {
	switch {
	case errors.Is(err, errs.ErrInvalidInput):
		return *errs.NewHTTPError(http.StatusBadRequest, err.Error())
	case errors.Is(err, errs.ErrinternalError):
		return *errs.NewHTTPError(http.StatusInternalServerError, err.Error())
	case errors.Is(err, errs.ErrNotFound):
		return *errs.NewHTTPError(http.StatusNotFound, err.Error())
	default:
		return *errs.NewHTTPError(http.StatusInternalServerError, "internal error")
	}
}
