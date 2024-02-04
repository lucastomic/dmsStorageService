package environment

import "os"

// GetProjectRoot returns the root directory of the project as specified by the "PROJECT_ROOT" environment variable.
// It is a utility function used to retrieve the base path for project-related files and operations.
func GetProjectRoot() string {
	return os.Getenv("PROJECT_ROOT")
}
