package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/lucastomic/dmsStorageService/internal/contextypes"
	"github.com/lucastomic/dmsStorageService/internal/controller/apitypes"
	"github.com/lucastomic/dmsStorageService/internal/errs"
	"github.com/lucastomic/dmsStorageService/internal/fileutils"
	"github.com/lucastomic/dmsStorageService/internal/logging"
	"github.com/lucastomic/dmsStorageService/internal/storageservice"
)

// StorageController manages file upload and retrieve operations from the local storage.
// It leverages a storage service for handling file storage and a common controller
// for shared HTTP handling logic.
type StorageController struct {
	logger         logging.Logger
	storageservice storageservice.StorageService
	common         CommonController
}

// New creates a new instance of StorageController with the provided logger and storage service.
// It initializes a StorageController that is responsible for handling file upload requests.
func New(
	logger logging.Logger,
	storageservice storageservice.StorageService,
) Controller {
	return &StorageController{logger, storageservice, CommonController{}}
}

// Router defines the routes that the StorageController handles.
// Currently, it sets up a single route for file upload requests, using the HTTP POST method.
func (c *StorageController) Router() apitypes.Router {
	return []apitypes.Route{
		{
			Path:    "/file",
			Method:  "POST",
			Handler: c.Upload,
		},
		{
			Path:    "/file/{id}",
			Method:  "GET",
			Handler: c.Get,
		},
	}
}

// Upload handles file upload requests.
// It validates the request, processes the file upload through the storage service,
// and returns an appropriate HTTP response.
func (c *StorageController) Upload(
	w http.ResponseWriter,
	req *http.Request,
) apitypes.Response {
	uploadData, err := c.parseAndValidateUploadReq(req, w)
	if err != nil {
		return c.common.ParseError(req.Context(), req, w, err)
	}
	err = c.storageservice.Upload(req.Context(), uploadData)
	if err != nil {
		return c.common.ParseError(req.Context(), req, w, err)
	}
	return apitypes.Response{
		Status:  http.StatusCreated,
		Content: "File stored successfully",
		Headers: map[string]string{"Content-Type": "application/json"},
	}
}

// Get handles the retrieval of a file based on its ID from the request's path variable.
// It validates the presence and type of the ID, retrieves the file, and returns it in the HTTP response.
// On failure, it constructs an appropriate error response.
func (c *StorageController) Get(w http.ResponseWriter, req *http.Request) apitypes.Response {
	id, err := extractIDFromRequest(req)
	if err != nil {
		return apitypes.Response{
			Status:  http.StatusBadRequest,
			Content: map[string]string{"error": err.Error()},
		}
	}
	file, err := c.storageservice.Get(req.Context(), id)
	if err != nil {
		return c.common.ParseError(req.Context(), req, w, err)
	}
	contentType, err := fileutils.DetermineMIME(file)
	if err != nil {
		return c.common.ParseError(req.Context(), req, w, err)
	}
	return apitypes.Response{
		Status: http.StatusOK,
		Headers: map[string]string{
			"Content-Disposition": fmt.Sprintf("attachment; filename=%s", file.Name()),
			"Content-Type":        contentType,
		},
		Content: file,
	}
}

// parseAndValidateUploadReq parses the incoming HTTP request to validate and extract necessary information
// for the file upload, such as ensuring the file size is within limits and extracting file metadata.
// It returns structured upload data or an error if validation fails.
// It checks the request is not over the maximum size and the form values Id and uploadfile exist.
func (c *StorageController) parseAndValidateUploadReq(
	req *http.Request,
	w http.ResponseWriter,
) (storageservice.UploadData, error) {
	const maxUploadSize = 10 << 20 // 10MB
	req.Body = http.MaxBytesReader(w, req.Body, maxUploadSize)

	if err := req.ParseMultipartForm(maxUploadSize); err != nil {
		return storageservice.UploadData{}, fmt.Errorf(
			"%w:%s",
			errs.ErrInvalidInput,
			"The uploaded file is too big. Maximum file size is 10MB.",
		)
	}

	id, err := strconv.ParseInt(req.FormValue("Id"), 10, 64)
	if err != nil {
		return storageservice.UploadData{}, fmt.Errorf(
			"%w:%s",
			errs.ErrInvalidInput,
			"Id must be an integer.",
		)
	}

	file, fileHeader, err := req.FormFile("uploadFile")
	if err != nil {
		return storageservice.UploadData{}, fmt.Errorf(
			"%w:%s",
			errs.ErrInvalidInput,
			"could not read uploaded file.",
		)
	}
	defer file.Close()

	return storageservice.UploadData{
		File:     file,
		Filename: fileHeader.Filename,
		Id:       id,
	}, nil
}

// extractIDFromRequest extracts and parses the ID path variable from the request context.
// Returns the parsed ID as int64 or an error if the ID is missing or not an integer.
func extractIDFromRequest(req *http.Request) (int64, error) {
	idParam := req.Context().Value(contextypes.ContextPathVarKey("id"))
	if idParam == nil {
		return 0, errors.New("id param can't be null")
	}

	idString, ok := idParam.(string)
	if !ok {
		return 0, errors.New("id param type is invalid")
	}

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return 0, errors.New("id param type is invalid")
	}
	return id, nil
}
