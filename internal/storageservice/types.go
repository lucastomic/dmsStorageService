package storageservice

import "mime/multipart"

// UploadData encapsulates the data required to upload a file.
// It contains the file stream, the name of the file, and an identifier
// that can be used to reference the file within the storage system.
type UploadData struct {
	File     multipart.File // File is the file stream to be uploaded.
	Filename string         // Filename is the name of the file.
	Id       int64          // Id is a unique identifier for the file.
}
