package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test POST /application with valid payload
func TestCreateApplication(t *testing.T) {
	router := setupRouter()
	server := httptest.NewServer(router)
	defer server.Close()

	testCases := []struct {
		description        string
		filePath           string
		expectedAssertions func(response *http.Response) bool
		expectedError      error
	}{
		{
			description: "Valid Payload",
			filePath:    "./samples/validPayloads/validPayload1.yaml",
			expectedAssertions: func(response *http.Response) bool {
				return assert.Equal(t, http.StatusOK, response.StatusCode)
			},
			expectedError: nil,
		},
		{
			description: "Valid Payload 2",
			filePath:    "./samples/validPayloads/validPayload2.yaml",
			expectedAssertions: func(response *http.Response) bool {
				return assert.Equal(t, http.StatusOK, response.StatusCode)
			},
			expectedError: nil,
		},
		{
			description: "Valid Payload 3",
			filePath:    "./samples/validPayloads/validPayload3.yaml",
			expectedAssertions: func(response *http.Response) bool {
				return assert.Equal(t, http.StatusOK, response.StatusCode)
			},
			expectedError: nil,
		},
		{
			description: "Valid Payload 4",
			filePath:    "./samples/validPayloads/validPayload4.yaml",
			expectedAssertions: func(response *http.Response) bool {
				return assert.Equal(t, http.StatusOK, response.StatusCode)
			},
			expectedError: nil,
		},
		{
			description: "Invalid Payload - missing title",
			filePath:    "./samples/invalidPayloads/missingTitle.yaml",
			expectedAssertions: func(response *http.Response) bool {
				return assert.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			expectedError: errors.New("title is required"),
		},
		{
			description: "Invalid Payload - missing version",
			filePath:    "./samples/invalidPayloads/missingVersion.yaml",
			expectedAssertions: func(response *http.Response) bool {
				return assert.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			expectedError: errors.New("version is required"),
		},
		{
			description: "Invalid Payload - missing maintainers",
			filePath:    "./samples/invalidPayloads/missingMaintainer.yaml",
			expectedAssertions: func(response *http.Response) bool {
				return assert.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			expectedError: errors.New("at least one maintainer is required"),
		},
		{
			description: "Invalid Payload - invalid maintainer email",
			filePath:    "./samples/invalidPayloads/invalidEmail.yaml",
			expectedAssertions: func(response *http.Response) bool {
				return assert.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			expectedError: errors.New("invalid maintainer email"),
		},
		{
			description: "Invalid Payload - missing company",
			filePath:    "./samples/invalidPayloads/missingCompany.yaml",
			expectedAssertions: func(response *http.Response) bool {
				return assert.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			expectedError: errors.New("company is required"),
		},
		{
			description: "Invalid Payload - missing website",
			filePath:    "./samples/invalidPayloads/missingWebsite.yaml",
			expectedAssertions: func(response *http.Response) bool {
				return assert.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			expectedError: errors.New("website is required"),
		},
		{
			description: "Invalid Payload - missing source",
			filePath:    "./samples/invalidPayloads/missingSource.yaml",
			expectedAssertions: func(response *http.Response) bool {
				return assert.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			expectedError: errors.New("source is required"),
		},
		{
			description: "Invalid Payload - missing license",
			filePath:    "./samples/invalidPayloads/missingLicense.yaml",
			expectedAssertions: func(response *http.Response) bool {
				return assert.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			expectedError: errors.New("license is required"),
		},
		{
			description: "Invalid Payload - missing description",
			filePath:    "./samples/invalidPayloads/missingDescription.yaml",
			expectedAssertions: func(response *http.Response) bool {
				return assert.Equal(t, http.StatusBadRequest, response.StatusCode)
			},
			expectedError: errors.New("description is required"),
		},
	}

	for _, tc := range testCases {
		t.Logf(tc.description)
		// Read the request payload from the file
		body, err := os.ReadFile(tc.filePath)
		assert.NoError(t, err)

		response := sendRequest(t, server, "POST", "/applications", body)
		assert.Equal(t, tc.expectedAssertions(response), true, "Expected assertions should pass")
	}
}

// Test GET /applications
func TestGetApplications(t *testing.T) {
	router := setupRouter()

	server := httptest.NewServer(router)
	defer server.Close()

	testCases := []struct {
		description        string
		reqPath            string
		expectedAssertions func(response *http.Response) bool
	}{
		{
			description: "List all applications",
			reqPath:     "/applications",
			expectedAssertions: func(response *http.Response) bool {
				var applications []Application
				err := json.NewDecoder(response.Body).Decode(&applications)

				return assert.Equal(t, http.StatusOK, response.StatusCode) && assert.NoError(t, err) &&
					assert.Equal(t, len(applications), 4)
			},
		},
		{
			description: "List all applications",
			reqPath:     "/applications?title=App",
			expectedAssertions: func(response *http.Response) bool {
				var applications []Application
				err := json.NewDecoder(response.Body).Decode(&applications)

				return assert.Equal(t, http.StatusOK, response.StatusCode) && assert.NoError(t, err) &&
					assert.Equal(t, len(applications), 4)
			},
		},
		{
			description: "List all applications",
			reqPath:     "/applications?title=App&version=1.0.0",
			expectedAssertions: func(response *http.Response) bool {
				var applications []Application
				err := json.NewDecoder(response.Body).Decode(&applications)

				return assert.Equal(t, http.StatusOK, response.StatusCode) && assert.NoError(t, err) &&
					assert.Equal(t, len(applications), 2)
			},
		},
	}

	for _, tc := range testCases {
		t.Logf(tc.description)

		response := sendRequest(t, server, "GET", tc.reqPath, nil)
		defer response.Body.Close()

		assert.Equal(t, tc.expectedAssertions(response), true, "Expected assertions should pass")
	}
}

func TestGetApplication(t *testing.T) {
	router := setupRouter()

	server := httptest.NewServer(router)
	defer server.Close()

	testCases := []struct {
		description        string
		reqPath            string
		expectedAssertions func(response *http.Response) bool
	}{
		{
			description: "Valid GET",
			reqPath:     "/applications/1",
			expectedAssertions: func(response *http.Response) bool {
				return assert.Equal(t, http.StatusOK, response.StatusCode)
			},
		}, {
			description: "Invalid GET",
			reqPath:     "/applications/23",
			expectedAssertions: func(response *http.Response) bool {
				return assert.Equal(t, http.StatusNotFound, response.StatusCode)
			},
		},
	}

	for _, tc := range testCases {
		t.Logf(tc.description)

		response := sendRequest(t, server, "GET", tc.reqPath, nil)
		defer response.Body.Close()

		assert.Equal(t, tc.expectedAssertions(response), true, "Expected assertions should pass")
	}
}

func TestDeleteApplication(t *testing.T) {
	router := setupRouter()

	server := httptest.NewServer(router)
	defer server.Close()

	testCases := []struct {
		description        string
		reqPath            string
		expectedAssertions func(response *http.Response) bool
	}{
		{
			description: "Valid DELETE",
			reqPath:     "/applications/1",
			expectedAssertions: func(response *http.Response) bool {
				var applications []Application
				err := json.NewDecoder(response.Body).Decode(&applications)

				return assert.Equal(t, http.StatusOK, response.StatusCode) && assert.NoError(t, err) &&
					assert.Equal(t, len(applications), 3)
			},
		},
	}

	for _, tc := range testCases {
		t.Logf(tc.description)

		response := sendRequest(t, server, "DELETE", tc.reqPath, nil)
		defer response.Body.Close()

		assert.Equal(t, tc.expectedAssertions(response), true, "Expected assertions should pass")
	}
}

// Helper function to send an HTTP request and return the response
func sendRequest(t *testing.T, server *httptest.Server, method, path string, body []byte) *http.Response {
	req, err := http.NewRequest(method, server.URL+path, bytes.NewReader(body))
	assert.NoError(t, err)

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/x-yaml")
	client := &http.Client{}
	response, err := client.Do(req)
	assert.NoError(t, err)

	return response
}
