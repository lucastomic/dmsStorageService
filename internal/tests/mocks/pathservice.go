package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockPathService struct {
	mock.Mock
}

func (m *MockPathService) Exists(ctx context.Context, id int64) (bool, error) {
	args := m.Called(ctx, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockPathService) GetPath(ctx context.Context, id int64) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

func (m *MockPathService) SavePath(ctx context.Context, id int64, path string) error {
	args := m.Called(ctx, id, path)
	return args.Error(0)
}
