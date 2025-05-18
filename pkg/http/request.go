package http

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    io.Reader
}

func NewRequest(method, url string) *Request {
	return &Request{
		Method:  method,
		URL:     url,
		Headers: make(map[string]string),
	}
}

func (r *Request) SetHeader(key, value string) *Request {
	r.Headers[key] = value
	return r
}

func (r *Request) SetBody(body []byte) *Request {
	r.Body = bytes.NewReader(body)
	return r
}

func (r *Request) SetBodyReader(body io.Reader) *Request {
	r.Body = body
	return r
}

func (r *Request) Build() (*http.Request, error) {
	parsedURL, err := url.Parse(r.URL)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(r.Method, parsedURL.String(), r.Body)
	if err != nil {
		return nil, err
	}

	// Set headers
	for key, value := range r.Headers {
		req.Header.Set(key, value)
	}

	return req, nil
}
