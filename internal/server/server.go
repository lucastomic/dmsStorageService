package server

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/lucastomic/dmsStorageService/internal/controller"
	"github.com/lucastomic/dmsStorageService/internal/controller/apitypes"
	"github.com/lucastomic/dmsStorageService/internal/logging"
	"github.com/lucastomic/dmsStorageService/internal/middleware"
)

// Server represents the core structure of an HTTP server. It encapsulates all necessary components
// for server operation, including routing, logging, and middleware management.
type Server struct {
	listenAddr  string                  // The address on which the server listens for incoming requests.
	controller  controller.Controller   // Controller manages routing of requests to their respective handlers.
	apiLogger   logging.Logger          // apiLogger is used for logging general API requests information.
	logicLogger logging.Logger          // logicLogger is specialized for logging business logic related events.
	middlewares []middleware.Middleware // middlewares is a slice of Middleware interfaces to be applied to all requests.
}

// New creates a new instance of the Server struct, initializing it with the provided parameters
// such as listen address, controller, API and logic loggers, and middlewares.
func New(
	listenAddr string,
	controller controller.Controller,
	apilogger logging.Logger,
	logicLogger logging.Logger,
	middlewares []middleware.Middleware,
) Server {
	return Server{
		listenAddr,
		controller,
		apilogger,
		logicLogger,
		middlewares,
	}
}

// Run initializes the server's routes based on the controller's router, applies middlewares,
// starts listening on the specified address, and logs the server's start or any errors encountered.
func (s *Server) Run() {
	r := mux.NewRouter()
	for _, route := range s.controller.Router() {
		handlerWithMiddlewares := middleware.ChainMiddleware(
			s.makeHTTPHandlerFunc(route.Handler),
			s.handleError,
			s.middlewares...,
		)
		r.Handle(route.Path, handlerWithMiddlewares).Methods(route.Method)
	}
	r.NotFoundHandler = middleware.ChainMiddleware(
		s.notFoundHandler,
		s.handleError,
		s.middlewares...,
	)
	s.apiLogger.Info(context.Background(), "Service running in %s", s.listenAddr)
	if err := http.ListenAndServe(s.listenAddr, r); err != nil {
		s.apiLogger.Error(context.Background(), "Failed to start server: %v", err)
	}
}

// makeHTTPHandlerFunc wraps the API function into an http.HandlerFunc, facilitating the handling
// of HTTP requests and responses within the server's routing mechanism.
func (s *Server) makeHTTPHandlerFunc(apiFn apitypes.APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := apiFn(w, r)
		s.writeResponse(r, w, res)
	}
}

// handleError handles errors by writing an error message as a JSON response, utilizing the
// writeJSON method to ensure the response format is consistent.
func (s *Server) handleError(
	req *http.Request,
	w http.ResponseWriter,
	err error,
	statusCode int,
) {
	errMsg := map[string]string{"error": err.Error()}
	s.writeResponse(req, w, apitypes.Response{Status: statusCode, Content: errMsg})
}

// notFoundHandler handles the request that points to an unexistent path.
// It returns a 404 error not found and prints that the page wan not found.
func (s *Server) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	s.writeResponse(r, w, apitypes.Response{Status: http.StatusNotFound, Content: "Page not found"})
}

// writeResponse prepares and sends an HTTP response based on the provided apitypes.Response struct.
// It sets custom headers, writes the status code, and sends the response content, which can vary in type.
// For *os.File types, the file's content is streamed to the response. For other types, the content is
// JSON-encoded and written to the response. Errors during JSON encoding are logged.
func (s *Server) writeResponse(req *http.Request, w http.ResponseWriter, res apitypes.Response) {
	setCustomHeaders(w, res.Headers)
	w.WriteHeader(res.Status)
	if err := writeContent(w, res.Content); err != nil {
		s.logicLogger.Error(req.Context(), "Failed to write response: %v", err)
	}
}

// writeContent writes the content to the response writer based on the content type.
// It handles *os.File by streaming the file's content and other types by JSON-encoding them.
// Returns an error if it encounters an issue during the write operation.
func writeContent(w http.ResponseWriter, content interface{}) error {
	switch c := content.(type) {
	case *os.File:
		defer c.Close()
		_, err := io.Copy(w, c)
		return err
	default:
		return json.NewEncoder(w).Encode(content)
	}
}

// setCustomHeaders sets the headers provided in the response struct.
func setCustomHeaders(w http.ResponseWriter, headers map[string]string) {
	for key, value := range headers {
		w.Header().Set(key, value)
	}
}
