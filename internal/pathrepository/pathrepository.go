package pathrepository

import "context"

// PathRepository defines the interface for operations on path storage.
// It outlines methods for saving, retrieving, and checking the existence of paths associated with unique identifiers.
type PathRepository interface {
	// SavePath persists a path associated with a given id in the storage.
	// If the ID already exists SavePath would override it.
	// The operation might fail due to storage errors, in which case an error will be returned.
	SavePath(ctx context.Context, id int64, path string) error

	// GetPath retrieves the path associated with the given id from the storage.
	// If the id does not exist, it returns an empty string and a resource-not-found error.
	// For storage-related errors, an error is returned.
	GetPath(ctx context.Context, id int64) (string, error)

	// Exists checks whether a path associated with the given id exists in the storage.
	// It returns true if the path exists, false otherwise. Errors are returned for storage-related issues.
	Exists(ctx context.Context, id int64) (bool, error)
}
