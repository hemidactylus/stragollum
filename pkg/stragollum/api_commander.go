package stragollum

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// DataAPICommander is a helper for making HTTP POST requests to the Data API.
type DataAPICommander struct {
	url   string
	token *string
}

// NewDataAPICommander creates a new DataAPICommander with the given URL and optional token.
func NewDataAPICommander(url string, token *string) *DataAPICommander {
	return &DataAPICommander{
		url:   url,
		token: token,
	}
}

// URL returns the commander's URL.
func (c *DataAPICommander) URL() string {
	return c.url
}

// Token returns the commander's token (may be nil).
func (c *DataAPICommander) Token() *string {
	return c.token
}

// RawRequest sends a POST request with the given payload to the commander's URL.
// It sets headers from the provided map, if any.
// If a token is present in the DataAPICommander, it adds a "Token" header.
// It returns the response body as bytes and an error if any occurred (including non-2xx HTTP status codes).
func (ac *DataAPICommander) RawRequest(payload []byte, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("POST", ac.url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers from the map
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Add token header if available
	if ac.token != nil {
		req.Header.Set("Token", *ac.token)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return bodyBytes, nil
}

// Request sends a JSON request and parses the JSON response.
// It automatically sets the "Content-Type" and "Accept" headers to "application/json".
// Input and output are automatically marshalled/unmarshalled as JSON.
func (ac *DataAPICommander) Request(requestObj interface{}, responseObj interface{}) error {
	// Marshal request object to JSON
	payload, err := json.Marshal(requestObj)
	if err != nil {
		return fmt.Errorf("failed to marshal request to JSON: %w", err)
	}

	// Set JSON content type header
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	// Send raw request
	respBody, err := ac.RawRequest(payload, headers)
	if err != nil {
		return err
	}

	// Unmarshal JSON response
	if err := json.Unmarshal(respBody, responseObj); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}
