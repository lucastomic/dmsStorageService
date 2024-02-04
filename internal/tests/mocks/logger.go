package mocks

import (
	"context"
	"net/http"
	"time"

	"github.com/lucastomic/dmsStorageService/internal/logging"
)

func NewLoggerMock() logging.Logger {
	return MockLogger{}
}

type MockLogger struct{}

func (m MockLogger) Request(
	ctx context.Context,
	r *http.Request,
	statusCode int,
	duration time.Duration,
) {
}
func (m MockLogger) Info(ctx context.Context, format string, a ...interface{})  {}
func (m MockLogger) Error(ctx context.Context, format string, a ...interface{}) {}
