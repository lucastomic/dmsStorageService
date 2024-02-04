package storageservice

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/lucastomic/dmsStorageService/internal/errs"
	"github.com/lucastomic/dmsStorageService/internal/tests/mocks"
	"github.com/stretchr/testify/mock"
)

var uploadTests = []struct {
	animal   UploadData
	expected error
}{
	{
		UploadData{Id: 1},
		errs.ErrInvalidInput,
	},
}

func TestUploadMetjod(t *testing.T) {
	logger := new(mocks.MockLogger)
	ctx := context.TODO()
	pathService := new(mocks.MockPathService)
	pathService.On("Exists", ctx, int64(1)).Return(true, nil)
	pathService.On("SavePath", ctx, 1, mock.AnythingOfType("string")).Return(nil)
	service := New(logger, pathService)

	for i, tt := range uploadTests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			got := service.Upload(ctx, tt.animal)
			if errors.Is(tt.expected, got) {
				t.Errorf("Expected: %v, got: %v", tt.expected, got)
			}
		})
	}
}
