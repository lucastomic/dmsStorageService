package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lucastomic/dmsStorageService/internal/controller/apitypes"
)

// TestWriteResponse checks if the writeResponse correctly sets headers and writes the response.
func TestWriteResponse(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "http://example.com/foo", nil)
	response := apitypes.Response{
		Status: http.StatusOK,
		Content: map[string]interface{}{
			"message": "success",
		},
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	srv := Server{}
	srv.writeResponse(req, w, response)
	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type header to be set to application/json")
	}
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code to be %v, got %v", http.StatusOK, w.Code)
	}
	expectedBody := `{"message":"success"}`
	body := w.Body.String()
	if body != expectedBody+"\n" { // json.Encoder adds a newline
		t.Errorf("Expected body to be %v, got %v", expectedBody, body)
	}
}

func TestSetCustomHeaders(t *testing.T) {
	w := httptest.NewRecorder()
	headers := map[string]string{
		"X-Custom-Header": "value",
	}
	setCustomHeaders(w, headers)

	if w.Header().Get("X-Custom-Header") != "value" {
		t.Errorf("Expected X-Custom-Header to be set to value")
	}
}
