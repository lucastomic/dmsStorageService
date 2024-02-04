package apitypes

import (
	"net/http"
)

// APIFunc is a type that represents a function signature for API handlers.
// It takes an http.ResponseWriter and an *http.Request as arguments and returns a Response struct.
// This design allows for a consistent API response structure across all handlers.
type APIFunc func(http.ResponseWriter, *http.Request) Response

// Router is a Routes slice. Each route contains the path, that method that will handle and the handler.
type Router []Route

// Response struct encapsulates the standard structure for API responses.
// It includes a Status field for the HTTP status code and a Message field for the response payload,
// which can be of any type to allow flexibility in response content.
type Response struct {
	Status  int               // HTTP status code to be returned with the response.
	Content any               // The payload of the response, allowing for flexible data types.
	Headers map[string]string // HTTP headers to be returned witht he response.
}

// Route is a struct with the necessary information for defaining and endpoint. Its path,
// method (POST, GET, PUT, etc.) and handler.
type Route struct {
	Path    string  // The route's path. E.g. /someroute/anotherone
	Method  string  // The HTTP method that will manage, like POST, PUT, GET, etc.
	Handler APIFunc // The function that will handle the route
}
