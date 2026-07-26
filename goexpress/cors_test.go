package goexpress

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSMiddleware(t *testing.T) {
    engine := NEW()
    engine.Use(CORS())

    var coreHandlerExecuted bool

    engine.GET("/test", func(c *Context) {
        coreHandlerExecuted = true
        c.String(http.StatusOK, "success")
    })

    // --- TEST 1: Standard GET Request ---
    reqGET := httptest.NewRequest("GET", "/test", nil)
    wGET := httptest.NewRecorder()
    engine.ServeHTTP(wGET, reqGET)

    // Verify Headers were injected
    if wGET.Header().Get("Access-Control-Allow-Origin") != "*" {
        t.Errorf("Expected CORS Origin header to be '*', got '%s'", wGET.Header().Get("Access-Control-Allow-Origin"))
    }
    
    // Verify core handler was reached
    if !coreHandlerExecuted {
        t.Errorf("Expected core handler to execute on GET request, but it didn't.")
    }

    // --- TEST 2: Preflight OPTIONS Request ---
    coreHandlerExecuted = false // Reset state
    engine.OPTIONS("/test", func(c *Context) {
        coreHandlerExecuted = true // This should NEVER execute due to c.Abort()
    })

    reqOPTIONS := httptest.NewRequest("OPTIONS", "/test", nil)
    wOPTIONS := httptest.NewRecorder()
    engine.ServeHTTP(wOPTIONS, reqOPTIONS)

    // Verify Status is 204 No Content
    if wOPTIONS.Code != http.StatusNoContent {
        t.Errorf("Expected preflight status 204, got %d", wOPTIONS.Code)
    }

    // Verify Abort() worked correctly
    if coreHandlerExecuted {
        t.Errorf("FATAL: Core handler executed during an OPTIONS preflight request! Abort() failed.")
    }
}
