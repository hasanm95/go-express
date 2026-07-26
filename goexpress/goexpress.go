package goexpress

import (
	"fmt"
	"net/http"
)

type RouteHandler func(c *Context)

type Engine struct {
	router *router
}

func NEW() *Engine{
	return &Engine{
		router: newRouter(),
	}
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
	node, params := e.router.getRoute(r.Method, r.URL.Path)

	if node != nil {
		c := newContext(w, r)
		c.Params = params
		key := r.Method + "-" + node.pattern
        e.router.handlers[key](c)
	} else {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL)
    }
}

func (e *Engine) Run (addr string) error {
	fmt.Printf("GoExpress is running on %s...\n", addr)
	return http.ListenAndServe(addr, e)
}