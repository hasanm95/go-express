package goexpress

import (
	"encoding/json"
	"net/http"
)

// Context wraps the standard library request and response,
// providing a simplified developer experience.
type Context struct {
	Writer http.ResponseWriter
	Req *http.Request
	Path string
	Method string
	Params map[string]string
}

// newContext is the internal factory for creating a Context per request
func newContext(w http.ResponseWriter, r *http.Request) *Context{
	return &Context{
		Writer: w,
		Req: r,
		Path: r.URL.Path,
		Method: r.Method,
	}
}


// Param retrieves a dynamic path parameter by its name
func (c *Context) Param(key string) string {
	return c.Params[key]
}


// --- HELPER METHODS ---

// String sends a plain text response with a status code
func (c *Context) String(code int, text string) {
	c.Writer.Header().Set("Content-Type", "text/plain")
	c.Writer.WriteHeader(code)
	c.Writer.Write([]byte(text))
}

// JSON sends a formatted JSON response with a status code
func (c *Context) JSON(code int, obj interface {}) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(code)

	// Use Go's built-in JSON encoder to translate the object and write it
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}  
}

// BindJSON reads the incoming HTTP request body and decodes it into Go struct
func (c *Context) BindJSON(obj interface{}) error{
	decoder := json.NewDecoder(c.Req.Body)
	defer c.Req.Body.Close()
	return decoder.Decode(obj)
}