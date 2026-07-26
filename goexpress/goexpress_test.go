package goexpress

import (
	"bytes"
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

// Test 4: Testing POST method widh BindJSON
func TestPostAndBindJSON(t *testing.T){
	engine := NEW()

	// 1. Define a dummy struct for testing
	type Payload struct {
		Message string `json:"message"`
	}

	// 2. Register a POST route that reads JSON and echoes it Back
	engine.POST("/echo", func(c *Context) {
		var p Payload
		if err := c.BindJSON(&p); err != nil {
			c.String(http.StatusBadRequest, "bad request")
			return
		}

		// If successful, send a 201 created status
		c.String(http.StatusCreated, "Received: "+p.Message)
	})

	// 3. Forge the incoming JSON body
	// bytes.NewBuffer turns a raw stream into io.Reader stream that httptest requires
	jsonBody := []byte(`{"message": "hello framework"}`)
	bodyReader := bytes.NewBuffer(jsonBody)

	// 4. Create a fake POST request
	req := httptest.NewRequest("POST", "/echo", bodyReader)
	w := httptest.NewRecorder()

	// 5. Execute
	engine.ServeHTTP(w, req)

	// 6. Assertions
	if w.Code != http.StatusCreated {
		t.Errorf("Expected sttus code 201, got %d", w.Code)
	}

	expectedBody := "Received: hello framework"

	if w.Body.String() != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, w.Body.String())
	}
}

// Test 5: Testing Put request
func TestPutRouting(t *testing.T) {
	engine := NEW()

	// 1. Define a dummy struct for testing
	type Payload struct {
		Role string `json:"role"`
	} 

	// 2. Register put routes for updates
	engine.PUT("/users/1/role", func(c *Context) {
		var p Payload

		if err := c.BindJSON(&p); err != nil {
			c.String(http.StatusBadRequest, "bad reqest")
			return
		}

		c.String(http.StatusOK, "Role updated to: "+p.Role)
	})

	// 3. Forge the incoming JSON body
	jsonBody := []byte(`{"Role": "Admin"}`)
	bodyReader := bytes.NewBuffer(jsonBody)

	// 4. Create fake PUT request
	req := httptest.NewRequest("PUT", "/users/1/role", bodyReader)
	w := httptest.NewRecorder()

	// 5. Execute
	engine.ServeHTTP(w, req)

	// 6. Assertions
	if w.Code != http.StatusOK {
		t.Errorf("Expected sttus code 200, got %d", w.Code)
	}

	expectedBody := "Role updated to: Admin"

	if w.Body.String() != expectedBody {
		t.Errorf("Expected body '%s', got '%s'", expectedBody, w.Body.String())
	}
}

// 6. Testing DELETE route
func TestDeleteRouting(t *testing.T) {
	engine := NEW()

	// 1. Setup mock route
	engine.DELETE("/remove", func(c *Context) {
		c.String(http.StatusOK, "Deleted")
	})

	// 2. Make fake request
	req := httptest.NewRequest("DELETE", "/remove", nil)
	w := httptest.NewRecorder()

	// 3. Execute
	engine.ServeHTTP(w, req)

	// 4. Assertions
	if w.Code != http.StatusOK {
		t.Errorf("Expected sttus code 200, got %d", w.Code)
	}
}

// 7. Testing dynamic route
func TestDynamicRoute(t *testing.T) {
	engine := NEW()

	// 1. Setup a mock route
	engine.GET("/users/:id", func(c *Context) {
		param := c.Param("id")

		c.String(http.StatusOK, "User Id is: "+param)
	})

	// 2. Make a fake request
	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()

	// 3. Execute
	engine.ServeHTTP(w, req)

	// 4. Assertions
	if w.Code != http.StatusOK {
		t.Errorf("Expected sttus code 200, got %d", w.Code)
	}

	expected := "User Id is: 123"

	if w.Body.String() != expected {
		t.Errorf("Expected body is '%s', got '%s'", expected, w.Body.String())
	}
}

// 8. Testing wild card
func TestWildCardRoute(t *testing.T) {
	engine := NEW()

	// 1. Setup a mock route
	engine.GET("/assets/*filepath", func(c *Context) {
		param := c.Param("filepath")

		c.String(http.StatusOK, "File: "+param)
	})

	// 2. Make a fake request
	req := httptest.NewRequest("GET", "/assets/css/main.css", nil)
	w := httptest.NewRecorder()

	// 3. Execute
	engine.ServeHTTP(w, req)

	// 4. Assertions
	expected := "File: css/main.css"
	   if w.Body.String() != expected {
        t.Errorf("Expected body '%s', got '%s'", expected, w.Body.String())
    }
}

// 9. Testing middleware order
func TestMiddlewareOrder(t *testing.T) {
	engine := NEW()

	// We will track the order of execution in this slice
	var executionOrder []string

	// Middleware A (Outer Layer)
	engine.Use(func(c *Context) {
		executionOrder = append(executionOrder, "A_Before")
		c.Next() // Suspend and go deeper
		executionOrder = append(executionOrder, "A_After")
	})

	// Middleware B (Inner Layer)
	engine.Use(func(c *Context) {
		executionOrder = append(executionOrder, "B_Before")
		c.Next() // Suspend and go to core handler
		executionOrder = append(executionOrder, "B_After")
	})

	// Core Route Handler
	engine.GET("/test", func(c *Context) {
		executionOrder = append(executionOrder, "Core_Handler")
		c.String(http.StatusOK, "OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)


	// Validate the Onion Model flow
	expectedOrder := []string{"A_Before", "B_Before", "Core_Handler", "B_After", "A_After"}
	
	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("Expected %d steps, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, v := range executionOrder {
		if v != expectedOrder[i] {
			t.Errorf("At index %d: expected %s, got %s", i, expectedOrder[i], v)
		}
	}
}