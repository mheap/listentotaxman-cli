package testutil

import (
	"bytes"
	"io"
	"net/http"
	"strings"
)

// MockRoundTripper implements http.RoundTripper for testing HTTP clients
type MockRoundTripper struct {
	Response     *http.Response
	ResponseBody string // Store body as string to recreate for each request
	Err          error
	RequestCount int
	LastRequest  *http.Request
	Requests     []*http.Request
}

// RoundTrip implements the http.RoundTripper interface
func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.RequestCount++
	m.LastRequest = req
	m.Requests = append(m.Requests, req)

	if m.Err != nil {
		return nil, m.Err
	}

	// Create a fresh response with a new body reader for each call
	resp := &http.Response{
		StatusCode: m.Response.StatusCode,
		Body:       io.NopCloser(strings.NewReader(m.ResponseBody)),
		Header:     m.Response.Header,
	}

	return resp, nil
}

// NewMockRoundTripperSuccess creates a mock that returns a successful response
func NewMockRoundTripperSuccess(statusCode int, body string) *MockRoundTripper {
	return &MockRoundTripper{
		Response: &http.Response{
			StatusCode: statusCode,
			Header:     make(http.Header),
		},
		ResponseBody: body,
	}
}

// NewMockRoundTripperError creates a mock that returns an HTTP error response
func NewMockRoundTripperError(statusCode int, body string) *MockRoundTripper {
	return &MockRoundTripper{
		Response: &http.Response{
			StatusCode: statusCode,
			Header:     make(http.Header),
		},
		ResponseBody: body,
	}
}

// NewMockRoundTripperNetworkError creates a mock that returns a network error
func NewMockRoundTripperNetworkError(err error) *MockRoundTripper {
	return &MockRoundTripper{
		Err: err,
	}
}

// GetRequestBody reads and returns the body of the last request
func (m *MockRoundTripper) GetRequestBody() (string, error) {
	if m.LastRequest == nil || m.LastRequest.Body == nil {
		return "", nil
	}

	body, err := io.ReadAll(m.LastRequest.Body)
	if err != nil {
		return "", err
	}

	// Reset the body so it can be read again if needed
	m.LastRequest.Body = io.NopCloser(bytes.NewBuffer(body))

	return string(body), nil
}
