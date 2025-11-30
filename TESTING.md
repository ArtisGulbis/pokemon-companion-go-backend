# Go Testing Guide

A comprehensive guide to testing in Go, specifically for this Pokemon API project.

## Table of Contents
1. [Testing Basics](#testing-basics)
2. [File Organization](#file-organization)
3. [Running Tests](#running-tests)
4. [Testing HTTP Handlers](#testing-http-handlers)
5. [Table-Driven Tests](#table-driven-tests)
6. [Testing Middleware](#testing-middleware)
7. [Best Practices](#best-practices)

---

## Testing Basics

### The `testing` Package

Go has built-in testing support through the `testing` package. No external frameworks needed!

### Test Function Requirements

```go
func TestXxx(t *testing.T) {
    // Test code here
}
```

**Rules:**
1. Function name MUST start with `Test`
2. Must take exactly one parameter: `*testing.T`
3. Must be in a `_test.go` file

### The `t *testing.T` Parameter

This gives you methods to control the test:

```go
t.Error("message")        // Report error, continue test
t.Errorf("format", args)  // Report formatted error, continue test
t.Fatal("message")        // Report error, stop test immediately
t.Fatalf("format", args)  // Report formatted error, stop test
t.Log("message")          // Log a message (only shown with -v flag)
t.Run("name", func)       // Run a subtest
```

**When to use Fatal vs Error:**
- Use `t.Fatal` when continuing doesn't make sense (e.g., failed to decode JSON)
- Use `t.Error` for assertions where you want to see all failures

---

## File Organization

Tests live next to the code they test:

```
your-project/
├── main.go
├── main_test.go          ← Tests for main.go
├── middleware/
│   ├── middleware.go
│   └── middleware_test.go ← Tests for middleware.go
└── model/
    └── model.go          ← No tests yet
```

**Package Declaration:**
```go
// In middleware_test.go
package middleware  // Same package as middleware.go

// This lets you test private functions and variables
// (those that start with lowercase)
```

---

## Running Tests

### Basic Commands

```bash
# Run tests in current directory
go test

# Run with verbose output (shows each test)
go test -v

# Run tests in all subdirectories
go test ./...

# Run a specific test
go test -run TestWriteError

# Run tests with pattern matching
go test -run TestFetch  # Runs TestFetchPokemon, TestFetchPokemonWithConfigurableURL

# Show code coverage
go test -cover

# Show detailed coverage
go test -coverprofile=coverage.out
go tool cover -html=coverage.out  # Opens in browser
```

### Understanding Output

```bash
=== RUN   TestWriteError
=== RUN   TestWriteError/Bad_Request_Error
--- PASS: TestWriteError/Bad_Request_Error (0.00s)
--- PASS: TestWriteError (0.00s)
PASS
ok  	github.com/your/package	0.290s
```

- `=== RUN` - Test started
- `--- PASS` - Test passed
- `--- FAIL` - Test failed
- `(0.00s)` - Time taken
- `ok` - All tests in package passed

---

## Testing HTTP Handlers

### Key Package: `net/http/httptest`

This package provides tools for testing HTTP code:

1. **`httptest.NewRecorder()`** - Fake ResponseWriter that captures output
2. **`httptest.NewRequest()`** - Creates a fake HTTP request
3. **`httptest.NewServer()`** - Creates a real test HTTP server

### Basic Handler Test Pattern

```go
func TestHandler(t *testing.T) {
    // 1. Create a fake request
    req := httptest.NewRequest("GET", "/path", nil)

    // 2. Create a recorder to capture the response
    w := httptest.NewRecorder()

    // 3. Call your handler
    yourHandler(w, req)

    // 4. Verify the response
    if w.Code != http.StatusOK {
        t.Errorf("Expected 200, got %d", w.Code)
    }

    // Check headers
    if ct := w.Header().Get("Content-Type"); ct != "application/json" {
        t.Errorf("Wrong Content-Type: %s", ct)
    }

    // Check body
    var response map[string]string
    json.NewDecoder(w.Body).Decode(&response)
    // ... verify response data
}
```

### Testing Handlers That Make HTTP Calls

When your handler calls external APIs, use a mock server:

```go
func TestFetchPokemon(t *testing.T) {
    // Create a mock API server
    mockAPI := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Verify the request
        if r.URL.Path != "/25" {
            t.Errorf("Wrong path: %s", r.URL.Path)
        }

        // Return mock data
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"id": 25, "name": "pikachu"}`))
    }))
    defer mockAPI.Close()  // Always close!

    // Use mockAPI.URL in your handler
    // ...
}
```

**Why use mock servers?**
- Tests run faster (no network)
- Tests are reliable (no external dependencies)
- Tests can simulate errors (network failures, timeouts, etc.)
- You control the exact responses

---

## Table-Driven Tests

This is the **idiomatic Go pattern** for testing multiple scenarios.

### Basic Pattern

```go
func TestSomething(t *testing.T) {
    // Define test cases
    tests := []struct {
        name     string  // Describes this test case
        input    string  // Input to test
        expected string  // Expected output
    }{
        {
            name:     "Empty string",
            input:    "",
            expected: "default",
        },
        {
            name:     "Normal input",
            input:    "hello",
            expected: "HELLO",
        },
    }

    // Loop through test cases
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Run the test with this case's inputs
            result := yourFunction(tt.input)

            // Verify the result
            if result != tt.expected {
                t.Errorf("Expected %s, got %s", tt.expected, result)
            }
        })
    }
}
```

### Why Use Table-Driven Tests?

1. **Easy to add test cases** - Just add another struct to the slice
2. **Clear structure** - Inputs and expected outputs are explicit
3. **Independent tests** - Each case runs separately (if one fails, others still run)
4. **Readable output** - Each case is named in the test output

### Example from This Project

See `/Users/artisgulbis/Development/react-native/go-backend/main_test.go` - `TestWriteError` tests three different error scenarios in one test function.

---

## Testing Middleware

Middleware is a function that wraps an `http.Handler`. Testing it requires verifying:

1. It correctly modifies the request/response
2. It calls the next handler when appropriate
3. It blocks requests when appropriate

### Basic Middleware Test Pattern

```go
func TestMiddleware(t *testing.T) {
    // Track if next handler was called
    nextCalled := false

    // Create a simple "next" handler
    nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        nextCalled = true
        w.WriteHeader(http.StatusOK)
    })

    // Wrap it with your middleware
    wrappedHandler := YourMiddleware(nextHandler)

    // Create test request
    req := httptest.NewRequest("GET", "/test", nil)
    w := httptest.NewRecorder()

    // Call the wrapped handler
    wrappedHandler.ServeHTTP(w, req)

    // Verify behavior
    if !nextCalled {
        t.Error("Next handler should have been called")
    }
}
```

### Example from This Project

See `/Users/artisgulbis/Development/react-native/go-backend/middleware/middleware_test.go`:
- `TestAuthorization` - Verifies API key checking
- `TestLogging` - Verifies request logging
- `TestCreateStack` - Verifies middleware are applied in correct order

---

## Best Practices

### 1. Use Meaningful Test Names

```go
// Good
func TestFetchPokemon_WhenNotFound_ReturnsError(t *testing.T)

// Avoid
func TestFunc1(t *testing.T)
```

### 2. Test One Thing Per Test

Each test should verify one behavior:

```go
// Good - separate tests
func TestValidation_EmptyInput(t *testing.T) { }
func TestValidation_InvalidFormat(t *testing.T) { }

// Avoid - testing too many things
func TestValidation(t *testing.T) {
    // tests empty input, invalid format, valid input, etc.
}
```

### 3. Use Table-Driven Tests for Multiple Scenarios

When testing the same function with different inputs, use table-driven tests.

### 4. Clean Up Resources

Always close test servers and files:

```go
mockServer := httptest.NewServer(handler)
defer mockServer.Close()  // Runs when test finishes
```

### 5. Make Tests Independent

Each test should be able to run in isolation:

```go
// Good - each test creates its own data
func TestSomething(t *testing.T) {
    data := createTestData()
    // ...
}

// Avoid - sharing state between tests
var sharedData = []string{"test"}
func TestA(t *testing.T) { sharedData = append(...) }
func TestB(t *testing.T) { // relies on sharedData from TestA }
```

### 6. Test Error Cases

Don't just test the happy path:

```go
tests := []struct{
    name string
    // ...
}{
    {name: "Success case"},
    {name: "Invalid input"},
    {name: "Network error"},
    {name: "Timeout"},
}
```

### 7. Use Subtests for Organization

```go
t.Run("group name", func(t *testing.T) {
    t.Run("subtest 1", func(t *testing.T) { })
    t.Run("subtest 2", func(t *testing.T) { })
})
```

### 8. Make Error Messages Helpful

```go
// Good - shows expected vs actual
t.Errorf("Expected status %d, got %d", expected, actual)

// Less helpful
t.Error("Wrong status")
```

---

## Common Patterns in This Project

### Testing Error Response Helper

```go
// See: main_test.go - TestWriteError
// Tests that writeError correctly formats error responses
```

### Testing Handler with External API

```go
// See: main_test.go - TestFetchPokemon
// Uses httptest.NewServer to mock PokeAPI
```

### Testing Authorization Middleware

```go
// See: middleware/middleware_test.go - TestAuthorization
// Verifies API key checking behavior
```

### Testing Logging Middleware

```go
// See: middleware/middleware_test.go - TestLogging
// Verifies request logging without breaking response flow
```

---

## Quick Reference

### Creating Test Requests

```go
// Simple GET request
req := httptest.NewRequest("GET", "/path", nil)

// POST with JSON body
body := strings.NewReader(`{"key": "value"}`)
req := httptest.NewRequest("POST", "/path", body)

// With headers
req.Header.Set("Content-Type", "application/json")
req.Header.Set("X-API-KEY", "secret")

// With path parameters (Go 1.22+)
req.SetPathValue("id", "123")
```

### Capturing Responses

```go
w := httptest.NewRecorder()
handler.ServeHTTP(w, req)

// Access response
statusCode := w.Code
headers := w.Header()
body := w.Body.String()  // or w.Body.Bytes()
```

### Decoding JSON Responses

```go
var response map[string]string
err := json.NewDecoder(w.Body).Decode(&response)
if err != nil {
    t.Fatalf("Failed to decode: %v", err)
}
```

---

## Next Steps

1. **Add more test cases** - Think about edge cases
2. **Increase coverage** - Run `go test -cover ./...` to see uncovered code
3. **Test the model package** - Add tests for any validation logic
4. **Integration tests** - Test the entire stack together
5. **Benchmark tests** - Measure performance (prefix with `Benchmark` instead of `Test`)

## Further Reading

- Official Go testing package: https://pkg.go.dev/testing
- httptest package: https://pkg.go.dev/net/http/httptest
- Table-driven tests: https://go.dev/wiki/TableDrivenTests
