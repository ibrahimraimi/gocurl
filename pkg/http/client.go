package http

import (
	"crypto/tls"
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
}

type ClientOptions struct {
	Timeout         time.Duration
	SkipTLSVerify   bool
	FollowRedirects bool
}

func DefaultClientOptions() ClientOptions {
	return ClientOptions{
		Timeout:         30 * time.Second,
		SkipTLSVerify:   false,
		FollowRedirects: true,
	}
}

// NewClient creates a new HTTP client with the given options
func NewClient(options ClientOptions) *Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: options.SkipTLSVerify,
		},
	}

	var redirectPolicy func(req *http.Request, via []*http.Request) error
	if !options.FollowRedirects {
		redirectPolicy = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	client := &http.Client{
		Timeout:       options.Timeout,
		Transport:     transport,
		CheckRedirect: redirectPolicy,
	}

	return &Client{
		client: client,
	}
}

// Do executes an HTTP request and returns the response
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}
