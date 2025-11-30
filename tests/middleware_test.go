package middleware

// Testing middleware in Go follows the same patterns
// Middleware is just a function that wraps an http.Handler
import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TESTING MIDDLEWARE: Authorization
// Middleware tests verify that:
// 1. The middleware correctly modifies requests/responses
// 2. The middleware calls the next handler when appropriate
// 3. The middleware blocks requests when appropriate
func TestAuthorization(t *testing.T) {
	tests := []struct {
		name           string
		apiKey         string // The X-API-KEY header value
		expectStatus   int    // Expected HTTP status code
		expectNextCall bool   // Should the next handler be called?
	}{
		{
			name:           "Valid API Key",
			apiKey:         "my-secret-key",
			expectStatus:   http.StatusOK, // Should pass through
			expectNextCall: true,
		},
		{
			name:           "Missing API Key",
			apiKey:         "",
			expectStatus:   http.StatusUnauthorized, // Should block
			expectNextCall: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Track whether the next handler was called
			nextCalled := false

			// Create a simple "next" handler that sets a flag when called
			// In real code, this would be your actual route handler
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Success"))
			})

			// Wrap the next handler with our Authorization middleware
			// This is the key: middleware returns a new handler
			wrappedHandler := Authorization(nextHandler)

			// Create a test request
			req := httptest.NewRequest("GET", "/test", nil)

			// Set the API key header (or don't, depending on test case)
			if tt.apiKey != "" {
				req.Header.Set("X-API-KEY", tt.apiKey)
			}

			// Create a response recorder
			w := httptest.NewRecorder()

			// Call the wrapped handler
			wrappedHandler.ServeHTTP(w, req)

			// Verify the status code
			if w.Code != tt.expectStatus {
				t.Errorf("Expected status %d, got %d", tt.expectStatus, w.Code)
			}

			// Verify whether next handler was called
			if nextCalled != tt.expectNextCall {
				t.Errorf("Expected nextCalled=%v, got %v", tt.expectNextCall, nextCalled)
			}

			// If unauthorized, verify the error response
			if tt.expectStatus == http.StatusUnauthorized {
				var response map[string]string
				if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
					t.Fatalf("Failed to decode error response: %v", err)
				}
				if response["error"] != "Unauthorized" {
					t.Errorf("Expected error 'Unauthorized', got '%s'", response["error"])
				}
			}
		})
	}
}

// TESTING MIDDLEWARE: Logging
// The Logging middleware is trickier because it captures the status code
// We need to verify it logs without errors and passes through correctly
func TestLogging(t *testing.T) {
	tests := []struct {
		name           string
		handlerStatus  int    // Status code the next handler will write
		handlerMessage string // Message the next handler will write
	}{
		{
			name:           "Successful Request",
			handlerStatus:  http.StatusOK,
			handlerMessage: "OK",
		},
		{
			name:           "Not Found Request",
			handlerStatus:  http.StatusNotFound,
			handlerMessage: "Not Found",
		},
		{
			name:           "Server Error",
			handlerStatus:  http.StatusInternalServerError,
			handlerMessage: "Internal Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a "next" handler that returns a specific status
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.handlerStatus)
				w.Write([]byte(tt.handlerMessage))
			})

			// Wrap with Logging middleware
			wrappedHandler := Logging(nextHandler)

			// Create test request
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			// Call the wrapped handler
			// This will trigger logging to stdout via log.Printf
			wrappedHandler.ServeHTTP(w, req)

			// Verify the status code was correctly passed through
			if w.Code != tt.handlerStatus {
				t.Errorf("Expected status %d, got %d", tt.handlerStatus, w.Code)
			}

			// Verify the response body was correctly passed through
			if w.Body.String() != tt.handlerMessage {
				t.Errorf("Expected body '%s', got '%s'", tt.handlerMessage, w.Body.String())
			}

			// Note: We can't easily test the log output without more complex setup
			// In production, you might use a custom logger that can be captured
			// For now, we verify the middleware doesn't break the request flow
		})
	}
}

// TESTING MIDDLEWARE STACK: CreateStack
// This tests that multiple middleware are applied in the correct order
func TestCreateStack(t *testing.T) {
	// Track the order in which middleware execute
	executionOrder := []string{}

	// Create test middleware that records execution order
	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "middleware1-before")
			next.ServeHTTP(w, r)
			executionOrder = append(executionOrder, "middleware1-after")
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			executionOrder = append(executionOrder, "middleware2-before")
			next.ServeHTTP(w, r)
			executionOrder = append(executionOrder, "middleware2-after")
		})
	}

	// Final handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		executionOrder = append(executionOrder, "handler")
		w.WriteHeader(http.StatusOK)
	})

	// Create stack and wrap handler
	stack := CreateStack(middleware1, middleware2)
	wrappedHandler := stack(handler)

	// Make request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w, req)

	// Verify execution order
	// CreateStack applies middleware in reverse, so:
	// middleware1 wraps middleware2 wraps handler
	// Execution: m1-before -> m2-before -> handler -> m2-after -> m1-after
	expectedOrder := []string{
		"middleware1-before",
		"middleware2-before",
		"handler",
		"middleware2-after",
		"middleware1-after",
	}

	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("Expected %d executions, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if executionOrder[i] != expected {
			t.Errorf("Position %d: expected '%s', got '%s'", i, expected, executionOrder[i])
		}
	}
}

// TESTING wrappedWriter
// This tests the custom ResponseWriter wrapper used by Logging
func TestWrappedWriter(t *testing.T) {
	// Create a ResponseRecorder (acts as base ResponseWriter)
	recorder := httptest.NewRecorder()

	// Create our wrapped writer
	wrapped := &wrappedWriter{
		ResponseWriter: recorder,
		statusCode:     http.StatusOK, // Default status
	}

	// Test 1: Verify default status code
	if wrapped.statusCode != http.StatusOK {
		t.Errorf("Expected default status 200, got %d", wrapped.statusCode)
	}

	// Test 2: Write a custom status code
	wrapped.WriteHeader(http.StatusNotFound)

	// Verify it captured the status code
	if wrapped.statusCode != http.StatusNotFound {
		t.Errorf("Expected status 404, got %d", wrapped.statusCode)
	}

	// Verify it also wrote to the underlying ResponseWriter
	if recorder.Code != http.StatusNotFound {
		t.Errorf("Expected underlying recorder to have status 404, got %d", recorder.Code)
	}
}
