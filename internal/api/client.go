package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/openlibing/openlibing-cli/internal/config"
)

// Client wraps http.Client with OpenLibing-specific behavior:
// auth injection, base URL resolution, and retry logic.
type Client struct {
	httpClient *http.Client
	baseURL    string
	auth       *config.OpenLibingAuth
	maxRetries int
}

// NewClient creates a new API client.
func NewClient(cfg *config.Config, auth *config.Auth) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:    strings.TrimRight(cfg.Endpoint, "/"),
		auth:       &auth.OpenLibing,
		maxRetries: 3,
	}
}

// Do executes an HTTP request with auth and retry logic.
func (c *Client) Do(method, endpoint string, queryParams map[string]string, headers map[string]string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + "/" + strings.TrimLeft(endpoint, "/")

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Set query parameters
	if len(queryParams) > 0 {
		q := req.URL.Query()
		for k, v := range queryParams {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	// Default headers to avoid WAF blocks
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; openlibing-cli)")
	req.Header.Set("Accept", "application/json, text/plain, */*")

	// Set headers from SPC (may override defaults)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Inject auth if configured
	if c.auth.Token != "" {
		req.Header.Set("Authorization", c.auth.TokenType+" "+c.auth.Token)
	}
	// Cookie-based auth (browser session)
	if c.auth.Cookie != "" {
		req.Header.Set("Cookie", c.auth.Cookie)
	}
	// CSRF token (required with cookie auth)
	if c.auth.CSRFToken != "" {
		req.Header.Set("csrf-token-open-li-bing", c.auth.CSRFToken)
		req.Header.Set("Referer", "https://www.openlibing.com/ops/dashboard/open-source-project")
		req.Header.Set("Origin", "https://www.openlibing.com")
	}

	// Retry loop
	var lastErr error
	for attempt := 0; attempt < c.maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s, 4s
			time.Sleep(time.Duration(1<<(attempt-1)) * time.Second)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue // network error, retry
		}

		// Only retry on 5xx server errors
		if resp.StatusCode >= 500 {
			resp.Body.Close()
			lastErr = fmt.Errorf("server error: %d", resp.StatusCode)
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", c.maxRetries, lastErr)
}

// Get is a convenience method for GET requests with query params.
func (c *Client) Get(endpoint string, queryParams map[string]string) (*http.Response, error) {
	return c.Do("GET", endpoint, queryParams, nil, nil)
}

// Post is a convenience method for POST requests.
func (c *Client) Post(endpoint string, body io.Reader) (*http.Response, error) {
	return c.Do("POST", endpoint, nil, nil, body)
}
