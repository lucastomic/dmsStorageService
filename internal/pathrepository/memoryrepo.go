package pathrepository

import (
	"context"
	"fmt"

	"github.com/lucastomic/dmsStorageService/internal/errs"
	"github.com/lucastomic/dmsStorageService/internal/logging"
)

// MemoryRepository returns an instance of PathRepository that operates in memory.
// It utilizes a map for storing paths associated with their IDs and leverages a logging.Logger
// for logging operations. This repository is intended for scenarios where persistence
// beyond the application lifecycle is not required.
func MemoryRepository(l logging.Logger) PathRepository {
	b := make(map[int64]string)
	return &memoryRepository{
		logger: l,
		buffer: &b,
	}
}

// memoryRepository implements the PathRepository interface, providing an in-memory storage solution
// for paths. It uses a map to associate paths with int64 IDs and supports operations to check existence,
// save, and retrieve paths.
type memoryRepository struct {
	logger logging.Logger    // logger for logging any errors or informational messages.
	buffer *map[int64]string // buffer is a map that stores paths associated with their IDs.
}

// Exists checks if a path associated with the given ID exists in the repository.
// It returns true if the path exists, false otherwise. This method does not generate errors.
func (m *memoryRepository) Exists(ctx context.Context, id int64) (bool, error) {
	_, exists := (*m.buffer)[id]
	return exists, nil
}

// SavePath stores a path associated with an ID in the repository.
// It overwrites any existing path associated with the ID.
func (m *memoryRepository) SavePath(ctx context.Context, id int64, path string) error {
	(*m.buffer)[id] = path
	return nil
}

// GetPath retrieves a path associated with the given ID from the repository.
// It returns an error if the path does not exist, logging the error before returning.
func (m memoryRepository) GetPath(ctx context.Context, id int64) (string, error) {
	path, exists := (*m.buffer)[id]
	if !exists {
		errMsg := fmt.Sprintf("path with id %d not found", id)
		return "", fmt.Errorf("%w: %s", errs.ErrNotFound, errMsg)
	}
	return path, nil
}
