package goexpress

// RouterGroup manages prefixed routes and their specific middlewares.
type RouterGroup struct {
	prefix      string
	middlewares []RouteHandler // Middlewares specifically for this group
	engine      *Engine        // Pointer to the main engine to access the router
}

// Group creates a new sub-group from the current group.
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		engine: engine,
	}
	
	newGroup.middlewares = make([]RouteHandler, len(group.middlewares))
	copy(newGroup.middlewares, group.middlewares)

	return newGroup
}

// Use adds middlewares ONLY to the current group.
func (group *RouterGroup) Use(middlewares ...RouteHandler) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// addRoute pre-calculates the complete execution chain before the server even starts.
func (group *RouterGroup) addRoute(method string, comp string, handler RouteHandler) {
	pattern := group.prefix + comp
	
	// Pre-allocate the exact size needed: Length of middlewares + 1 for the core handler
	handlers := make([]RouteHandler, len(group.middlewares), len(group.middlewares)+1)
	copy(handlers, group.middlewares)
	handlers = append(handlers, handler)

	// Register the fully calculated array into the Trie tree
	group.engine.router.addRoute(method, pattern, handlers)
}

// GET defines the method to add a GET request
func (group *RouterGroup) GET(pattern string, handler RouteHandler) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add a POST request
func (group *RouterGroup) POST(pattern string, handler RouteHandler) {
	group.addRoute("POST", pattern, handler)
}

// PUT defines the method to fully replace or create a resource
func (group *RouterGroup) PUT(pattern string, handler RouteHandler) {
	group.addRoute("PUT", pattern, handler)
}

// PATCH defines the method to partially modify a resource
func (group *RouterGroup) PATCH(pattern string, handler RouteHandler) {
	group.addRoute("PATCH", pattern, handler)
}

// DELETE defines the method to remove a resource
func (group *RouterGroup) DELETE(pattern string, handler RouteHandler) {
	group.addRoute("DELETE", pattern, handler)
}

// OPTIONS defines the method to handle CORS preflight requests
func (group *RouterGroup) OPTIONS(pattern string, handler RouteHandler) {
    group.addRoute("OPTIONS", pattern, handler)
}
