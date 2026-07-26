package goexpress

import (
	"net/http"
)

// CORS returns a middleware that handles Cross-Origin Resource Sharing.
func CORS() RouteHandler {
    return func(c *Context) {
        // 1. Inject the necessary CORS Headers into EVERY response
        c.SetHeader("Access-Control-Allow-Origin", "*")
        c.SetHeader("Access-Control-Allow-Credentials", "true")
        
        // Whitelist the HTTP methods your framework allows
        c.SetHeader("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
        
        // Whitelist standard headers + our custom X-Admin-Token from Lab 6
        c.SetHeader("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Admin-Token")

        // 2. Handle the Preflight Request
        if c.Method == "OPTIONS" {
            // A preflight request doesn't need a response body.
            // It just needs the headers (which we attached above) and a 204 status.
            c.Writer.WriteHeader(http.StatusNoContent)
            
            // CRITICAL: We use Abort() from Lab 6 so the execution chain stops here!
            // We do not want preflight requests bleeding into our core logic.
            c.Abort()
            return
        }

        // 3. If it's a normal request (GET, POST, etc.), continue down the chain
        c.Next()
    }
}
