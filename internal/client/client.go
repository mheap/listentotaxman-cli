package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mheap/listentotaxman-cli/internal/types"
)

const apiURL = "https://listentotaxman.com/ws/tax/index.js.php"

// Client represents the API client
type Client struct {
	httpClient *http.Client
}

// New creates a new API client
func New() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

// NewWithHTTPClient creates a new API client with a custom HTTP client (for testing)
func NewWithHTTPClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
	}
}

// CalculateTax calls the listentotaxman API and returns the tax calculation
func (c *Client) CalculateTax(req *types.TaxRequest) (*types.TaxResponse, error) {
	// Set fixed fields
	req.Response = "json"
	req.Time = "1"

	// Marshal request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Execute request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var taxResp types.TaxResponse
	if err := json.Unmarshal(body, &taxResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &taxResp, nil
}
