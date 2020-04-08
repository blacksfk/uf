package microframework

// Stores the underlying endpoint, route-wide middleware, and server
type Group struct {
	endpoint   string
	middleware []Middleware
	server     *Server
}

// Add middleware to be called for handlers following this method call for this group
func (g *Group) Middleware(nextRoutes ...Middleware) *Group {
	g.middleware = append(g.middleware, nextRoutes...)

	return g
}

// Bind this route to support GET requests, with methodOnly middleware only applied here
func (g *Group) Get(h Handler, methodOnly ...Middleware) *Group {
	g.server.Get(g.endpoint, h, append(g.middleware, methodOnly...)...)

	return g
}

// Bind this route to support POST requests, with methodOnly middleware only applied here
func (g *Group) Post(h Handler, methodOnly ...Middleware) *Group {
	g.server.Post(g.endpoint, h, append(g.middleware, methodOnly...)...)

	return g
}

// Bind this route to support PUT requests, with methodOnly middleware only applied here
func (g *Group) Put(h Handler, methodOnly ...Middleware) *Group {
	g.server.Put(g.endpoint, h, append(g.middleware, methodOnly...)...)

	return g
}

// Bind this route to support PATCH requests, with methodOnly middleware only applied here
func (g *Group) Patch(h Handler, methodOnly ...Middleware) *Group {
	g.server.Patch(g.endpoint, h, append(g.middleware, methodOnly...)...)

	return g
}

// Bind this route to support DELETE requests, with methodOnly middleware only applied here
func (g *Group) Delete(h Handler, methodOnly ...Middleware) *Group {
	g.server.Delete(g.endpoint, h, append(g.middleware, methodOnly...)...)

	return g
}
