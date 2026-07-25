package goexpress

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test 1: Does the Engine successfully find and execute a valid route?
func TestEngineRouting(t *testing.T){
	// 1. Setup a fresh engine
	engine := NEW()

	// 2. Register a test route
	engine.GET("/hello", func(c *Context) {
		c.String(http.StatusOK, "Hello world")
	})

	// 3. Create a fake "offline" request (Simulating a browser)
	req := httptest.NewRequest("GET", "/hello", nil)

	// 4. Create a fake ResponseWriter to capture the output
	w := httptest.NewRecorder()

	// 5. Execute our Engine's core router directly
	engine.ServeHTTP(w, req)

	// 6. Assertions: Did the engnie do the right thing?
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200 but got %d", w.Code)
	}

	if w.Body.String() != "Hello world" {
		t.Errorf("Expected status code 'Hello world' but got '%s'", w.Body.String())
	}
}

// Test 2: Does the engine correctly return a 404 when a route missing?

func TestEngineNotFond(t *testing.T){
	// 1. Setup an Empty engine (we are not registering any routes)
	engine := NEW()

	// 2. Create a fake request for a page that doesn't exist
	req := httptest.NewRequest("GET", "/missing-page", nil)
	w := httptest.NewRecorder()

	// 3. Execute
	engine.ServeHTTP(w, req)

	// 4. Assertions
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %d", w.Code)
	}
}

// Test 3: Testing the context with JSON Helper
func TestContextWithJSON(t *testing.T) {
	engine := NEW()

	// 1. Setup a route that sends a JSON response
	engine.GET("/api/data", func(c *Context) {
		c.JSON(http.StatusCreated, map[string]string{"status": "ok"})
	})

	// 2. Create our fakes and execute
	req := httptest.NewRequest("GET", "/api/data", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	// 3. Assertion A: Check the status code
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code 201, got %d", w.Code)
	}

	// 4. Assertion B: Check the http headers
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected content type 'application/json', got '%s'", contentType)
	}

	// 5. Assertion C: Parse and verify the exact JSON payload
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected JSON status 'ok', got '%s'", response["status"])
	}
}