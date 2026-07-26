package goexpress

import (
	"fmt"
	"net/http"
)

type RouteHandler func(c *Context)

type Engine struct {
	*RouterGroup 
	router *router
}

func NEW() *Engine{
	engine := &Engine{router: newRouter()}
	// The Engine initializes itself as the absolute Root Group ("")
	engine.RouterGroup = &RouterGroup{engine: engine}
	return engine
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)

	node, params := e.router.getRoute(r.Method, r.URL.Path)
	if node != nil {
		c.Params = params
		key := r.Method + "-" + node.pattern
		
		// Fetch the pre-calculated execution chain directly from the router
		c.handlers = e.router.handlers[key] 
	} else {
		// 404 Handler
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 NOT FOUND")
		})
	}

	// Kick off the execution chain
	c.Next()
}

func (e *Engine) Run (addr string) error {
	fmt.Printf("GoExpress is running on %s...\n", addr)
	return http.ListenAndServe(addr, e)
}