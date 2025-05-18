package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Response represents an HTTP response with additional metadata
type Response struct {
	StatusCode int
	Status     string
	Headers    http.Header
	Body       []byte
}

// NewResponse creates a Response from an http.Response
func NewResponse(resp *http.Response) (*Response, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    resp.Header,
		Body:       body,
	}, nil
}

// String returns a string representation of the response
func (r *Response) String() string {
	return fmt.Sprintf("%s\n%s", r.Status, string(r.Body))
}

// JSON returns the response body as a JSON object
func (r *Response) JSON() (interface{}, error) {
	var result interface{}
	err := json.Unmarshal(r.Body, &result)
	return result, err
}

// ToJSON returns the response as a JSON string
func (r *Response) ToJSON() (string, error) {
	data := map[string]interface{}{
		"status_code": r.StatusCode,
		"status":      r.Status,
		"headers":     r.Headers,
		"body":        string(r.Body),
	}

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
