package goexpress

import (
	"fmt"
	"net/http"
)

type RouteHandler func(w http.ResponseWriter, r *http.Request)

type Engine struct {
	router map[string]RouteHandler
}

func NEW() *Engine{
	return &Engine{
		router: make(map[string]RouteHandler),
	}
}

func (e *Engine) GET (path string, handler RouteHandler) {
	e.router["GET-"+path] = handler
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := r.Method + "-" + r.URL.Path

	if routeHandler, ok := e.router[key]; ok {
		routeHandler(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", r.URL)
	}
}

func (e *Engine) Run (addr string) error {
	fmt.Printf("GoExpress is running on %s...\n", addr)
	return http.ListenAndServe(addr, e)
}