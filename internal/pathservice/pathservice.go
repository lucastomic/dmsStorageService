package pathservice

import (
	"context"
	"fmt"

	"github.com/lucastomic/dmsStorageService/internal/errs"
	"github.com/lucastomic/dmsStorageService/internal/logging"
	"github.com/lucastomic/dmsStorageService/internal/pathrepository"
)

// PathService defines the interface for managing paths associated with IDs.
// It abstracts the operations for saving, retrieving, and checking the existence of paths
// to enable flexible implementation details, such as storage mechanisms.
type PathService interface {
	SavePath(
		ctx context.Context,
		id int64,
		path string,
	) error // Saves a path associated with an ID.
	GetPath(
		ctx context.Context,
		id int64,
	) (string, error) // Retrieves a path associated with an ID. Returns an error if the id doesn't exist.
	Exists(
		ctx context.Context,
		id int64,
	) (bool, error) // Checks if a path associated with an ID exists.
}

// pathService implements the PathService interface, providing methods to interact
// with path storage through a repository pattern. It utilizes a logger for error logging.
type pathService struct {
	logger logging.Logger                // Used for logging errors encountered during operations.
	repo   pathrepository.PathRepository // The repository responsible for actual storage operations.
}

// New initializes a new instance of pathService with the provided logger and repository.
// This function serves as a constructor for pathService, encapsulating the dependencies needed for path management.
func New(logger logging.Logger, repo pathrepository.PathRepository) PathService {
	return pathService{logger, repo}
}

// Exists checks if a path associated with the given ID exists in the storage.
// It logs an error and returns false along with the error if the repository encounters an issue.
func (p pathService) Exists(ctx context.Context, id int64) (bool, error) {
	exists, err := p.repo.Exists(ctx, id)
	if err != nil {
		p.logger.Error(ctx, "Error checking path existence: %s", err.Error())
		return false, err
	}
	return exists, nil
}

// SavePath persists a path associated with an ID.
// It's important to take into account that SavePath will override the path if the ID it's already in use.
// In order to not override any ID, it can be performed Exists before Save.
// If the repository fails to save the path, it logs the error and returns an internal error wrapped around the original error.
func (p pathService) SavePath(ctx context.Context, id int64, path string) error {
	err := p.repo.SavePath(ctx, id, path)
	if err != nil {
		p.logger.Error(ctx, "Error saving path: %s", err.Error())
		return err
	}
	return nil
}

// GetPath retrieves a path associated with the given ID from the repository.
// If the id tryed to retrieve doesn't exist,returns an error
// It logs and returns any error encountered during the retrieval process.
func (p pathService) GetPath(ctx context.Context, id int64) (string, error) {
	pathExists, err := p.Exists(ctx, id)
	if err != nil {
		p.logger.Error(ctx, "Error checking path with id %d: %s", id, err.Error())
		return "", err
	}
	if !pathExists {
		return "", fmt.Errorf("failed retrieving ID %d: %w", id, errs.ErrNotFound)
	}
	path, err := p.repo.GetPath(ctx, id)
	if err != nil {
		p.logger.Error(ctx, "Error retrieving path: %s", err.Error())
		return "", err
	}
	return path, nil
}
