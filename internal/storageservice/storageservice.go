package storageservice

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/lucastomic/dmsStorageService/internal/environment"
	"github.com/lucastomic/dmsStorageService/internal/errs"
	"github.com/lucastomic/dmsStorageService/internal/logging"
	"github.com/lucastomic/dmsStorageService/internal/pathservice"
)

// StorageService defines the interface for storage operations, including uploading and retrieving files.
// It abstracts the underlying storage mechanism, allowing for different implementations.
type StorageService interface {
	// Upload processes and stores the provided UploadData in the storage.
	// It returns an error if the upload fails due to validation issues or storage errors.
	Upload(context.Context, UploadData) error

	// Get retrieves the file identified by the specified identifier.
	// It returns an error if the retrieval fails or if the file does not exist.
	Get(context.Context, int64) (*os.File, error)
}

// New initializes a new instance of a StorageService with the provided logger and path service.
// This constructor function returns a storageService that uses the given pathService for path management
// and the logger for logging errors and information.
func New(logger logging.Logger, pathservice pathservice.PathService) StorageService {
	return &storageService{logger, pathservice}
}

// storageService implements the StorageService interface, providing methods for file upload and retrieval.
// It utilizes a path service for managing file paths and a logger for error logging.
type storageService struct {
	logger  logging.Logger          // Logger for logging operations and errors.
	pathsrv pathservice.PathService // Path service for managing file paths.
}

// Upload handles the storage of given UploadData.
// It checks if the path already exists to prevent duplicates, stores the file, and then saves the path.
// Errors during these operations are logged and may result in a rollback of the file storage.
func (s *storageService) Upload(
	ctx context.Context,
	data UploadData,
) error {
	alreadyExists, err := s.pathsrv.Exists(ctx, data.Id)
	if err != nil {
		s.logger.Error(ctx, "Failed to check if ID exists: %s", err.Error())
		return err
	}
	if alreadyExists {
		return fmt.Errorf("path with id %d already exists: %w", data.Id, errs.ErrInvalidInput)
	}
	dest, err := s.storeFile(ctx, data)
	if err != nil {
		s.logger.Error(ctx, "Failed to store file: %s", err.Error())
		return fmt.Errorf("error storing the file: %w", errs.ErrinternalError)
	}
	err = s.pathsrv.SavePath(ctx, data.Id, dest)
	if err != nil {
		// TODO: Rollback file storage
		s.logger.Error(ctx, "Failed to save path: %s", err.Error())
		return err
	}
	return nil
}

// Get retrieves the file associated with the given ID from the storage.
// It first fetches the file path using the path service. If the path cannot be retrieved
// or if the file does not exist at the retrieved path, it logs an error and returns an error.
// If the id deosn't exist it returns an ErrNotFound error.
// On successfully locating the file, it opens the file for reading and returns the file handle.
func (s *storageService) Get(ctx context.Context, id int64) (*os.File, error) {
	filePath, err := s.pathsrv.GetPath(ctx, id)
	if err != nil {
		return &os.File{}, err
	}

	if checkIfDoesntExist(filePath) {
		s.logger.Error(ctx, "File with ID %d not found in path %s", id, filePath)
		return &os.File{}, fmt.Errorf(
			"file with ID %d can't be reached at its path",
			id,
		)
	}
	file, err := os.Open(filePath)
	if err != nil {
		s.logger.Error(ctx, "Failed to open file %d: %s", id, err.Error())
		return &os.File{}, fmt.Errorf(
			"failed to open file with ID %d: %w",
			id,
			errs.ErrinternalError,
		)
	}

	return file, nil
}

// storeFile is a helper method for storing a file on the filesystem.
// It generates the full path for the file, creates it, and copies the content from UploadData.
// Returns the path of the stored file or an error if the operation fails.
func (s storageService) storeFile(ctx context.Context, data UploadData) (string, error) {
	root := environment.GetProjectRoot()
	dirPath := filepath.Join(root, "files")
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}
	dst, err := os.Create(filepath.Join(dirPath, data.Filename))
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()
	if _, err := io.Copy(dst, data.File); err != nil {
		return "", fmt.Errorf("failed to write to file: %w", err)
	}
	return dst.Name(), nil
}

func checkIfDoesntExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return os.IsNotExist(err)
}
