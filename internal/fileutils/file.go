package fileutils

import (
	"io"
	"net/http"
	"os"
)

// DetermineMIME reads the first 512 bytes of the provided file to determine its MIME type
// using the http.DetectContentType function. After reading, it seeks back to the beginning
// of the file to ensure that subsequent operations on the file start from the correct position.
// This function is particularly useful for dynamically determining the content type of a file
// when serving it to HTTP clients, and the MIME type is not predetermined.
func DetermineMIME(file *os.File) (string, error) {
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}
	contentType := http.DetectContentType(buffer[:n])

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	return contentType, nil
}
