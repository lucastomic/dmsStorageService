package pathrepository

import (
	"context"
	"errors"
	"testing"

	"github.com/lucastomic/dmsStorageService/internal/errs"
	"github.com/lucastomic/dmsStorageService/internal/tests/mocks"
)

var (
	logger = mocks.NewLoggerMock()
	repo   = MemoryRepository(logger)
	ctx    = context.TODO()
	id     = int64(1)
	path   = "test/path"
)

func TestSavePath(t *testing.T) {
	repo.SavePath(ctx, id, path)

	exists, err := repo.Exists(ctx, id)
	if err != nil || !exists {
		t.Errorf("Expected path to exist after saving, got exists=%v, err=%v", exists, err)
	}
}

func TestGetPath(t *testing.T) {
	repo.SavePath(ctx, id, path)

	got, err := repo.GetPath(ctx, id)
	if err != nil || got != path {
		t.Errorf("Expected to retrieve path '%s', got '%s', err=%v", path, got, err)
	}
}

func TestGetPathFail(t *testing.T) {
	repo.SavePath(ctx, id, path)
	_, err := repo.GetPath(ctx, id+1) // Using an ID that has not been saved
	if err == nil {
		t.Errorf("Expected an error for an unsaved path ID, got nil")
	} else if !errors.Is(err, errs.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %T", err)
	}
}

func TestOverwrite(t *testing.T) {
	repo.SavePath(ctx, id, path)
	newPath := "new/test/path"
	repo.SavePath(ctx, id, newPath)
	path, _ = repo.GetPath(ctx, id)
	if path != newPath {
		t.Errorf("Expected path to be overwritten with '%s', got '%s'", newPath, path)
	}
}
