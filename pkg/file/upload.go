package file

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// UploadFile represents a file to be uploaded
type UploadFile struct {
	FieldName string
	FilePath  string
}

// CreateMultipartRequest creates a multipart request with file uploads
func CreateMultipartRequest(url string, files []UploadFile, extraFields map[string]string) (*http.Request, error) {
	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)

	// Add files to the multipart writer
	for _, file := range files {
		fileContents, err := os.Open(file.FilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", file.FilePath, err)
		}
		defer fileContents.Close()

		// Create a form file field
		fieldWriter, err := multipartWriter.CreateFormFile(file.FieldName, filepath.Base(file.FilePath))
		if err != nil {
			return nil, fmt.Errorf("failed to create form file: %w", err)
		}

		// Copy the file contents to the form field
		_, err = io.Copy(fieldWriter, fileContents)
		if err != nil {
			return nil, fmt.Errorf("failed to copy file contents: %w", err)
		}
	}

	// Add extra form fields
	for key, value := range extraFields {
		if err := multipartWriter.WriteField(key, value); err != nil {
			return nil, fmt.Errorf("failed to add form field: %w", err)
		}
	}

	// Close the multipart writer
	if err := multipartWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set the content type
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	return req, nil
}
