package goexpress

import (
	"fmt"
	"net/http"
)

type RouteHandler func(c *Context)

type Engine struct {
	router *router
	middlewares []RouteHandler // Store global middleware functions
}

func NEW() *Engine{
	return &Engine{
		router: newRouter(),
	}
}

// Use adds global middlewares into framework instance
func (e *Engine) Use(middlewares ...RouteHandler) {
	e.middlewares = append(e.middlewares, middlewares...)
}

func (e *Engine) addRoute(method string, pattern string, handler RouteHandler) {
    e.router.addRoute(method, pattern, handler)
}

func (e *Engine) GET (pattern string, handler RouteHandler) {
	e.addRoute("GET", pattern, handler)
}

func (e *Engine) POST (pattern string, handler RouteHandler) {
	e.addRoute("POST", pattern, handler)
}

func (e *Engine) PUT (pattern string, handler RouteHandler) {
	e.addRoute("PUT", pattern, handler)
}

func (e *Engine) DELETE (pattern string, handler RouteHandler) {
	e.addRoute("DELETE", pattern, handler)
}


func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. Create the briefcase
	c := newContext(w, r) 

	// 2. Preload the gloabl middlewares into execution chain
	c.handlers = append(c.handlers, e.middlewares...)

	// 3. Find the specific route logic
	node, params := e.router.getRoute(r.Method, r.URL.Path)
	if node != nil {
		c.Params = params
		key := r.Method + "-" + node.pattern
		// 4. Append core handler to the end of the chain
		c.handlers = append(c.handlers, e.router.handlers[key])
	} else {
		// Append a 404 handler to the chain if route is missing
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND")
		})
	}

	// 5. Kick off the execution chain
	c.Next()
}

func (e *Engine) Run (addr string) error {
	fmt.Printf("GoExpress is running on %s...\n", addr)
	return http.ListenAndServe(addr, e)
}