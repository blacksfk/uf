/*
A tiny framework implementing routing (via vestigo), middleware support, and error handling.

Example:

import (
	uf "microframework"
	// ...
)

func main() {
	// create a new server
	server := uf.NewServer(":9001")

	// add routes to the server (supports GET, POST, PUT, PATCH, DELETE convenience methods)
	// middleware functions A, B, C will be called before the handler
	server.Get("/book", handler, middlewareA, middlewareB, middlewareC)

	// ...

	// add route groups to the server
	// middleware functions A, B will be called before each middleware and handler defined below
	server.Group("/author", middlewareA, middlewareB)

		// middlewareC will be called after A and B but only for GET requests on this route
		.Get(author.HandleGet, middlewareC)

		// middlewareD will be called after A and B, for the following route definitions
		.Middleware(middlewareD)
		.Post(author.HandlePost)
		.Put(author.HandlePut)

		// middlewareE will be called after A, B, and D but only for PATCH requests on this route
		.Patch(author.HandlePatch, middlewareE)
		.Delete(author.HandleDelete)

	// start the server
	e := server.Start()

	// ...
}

func handler(w http.ResponseWriter, r *http.Request) error {
	books := database.ObtainBooks()

	return uf.SendJSON(books)
}

func middlewareA(next uf.Handler) uf.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		// get the auth key and user (somehow)
		key := r.Header.Get("Authorization")
		user := database.FindUser()

		if !user.Valid(key) {
			// user needs to re-authenticate
			return uf.Unauthorized("Invalid login")
		}

		// authenticated
		r = r.WithContext(user.ToContext(r.Context()))

		// progress to next handler
		return next(w, r)
	}
}
*/
package microframework

import (
	"fmt"
	"github.com/husobee/vestigo"
	"net/http"
)

// Extension of http.Handler that returns an error to the framework's error handler
type Handler func(http.ResponseWriter, *http.Request) error

// Middleware are setup functions called once during startup that return intermediary
// Handler functions which are called after route matching but before the controller
type Middleware func(Handler) Handler

// Wrapper around vestigo.Router
type Server struct {
	address string
	*vestigo.Router
}

// create a new server
func NewServer(address string) *Server {
	router := vestigo.NewRouter()

	return &Server{address, router}
}

// listen for connections on the previously supplied address
func (server *Server) Start() error {
	fmt.Printf("Starting server on %s\n", server.address)

	return http.ListenAndServe(server.address, server)
}

func (server *Server) Get(endpoint string, c Handler, m ...Middleware) {
	server.Add(http.MethodGet, endpoint, newQueue(c, m).ServeHTTP)
}

func (server *Server) Post(endpoint string, c Handler, m ...Middleware) {
	server.Add(http.MethodPost, endpoint, newQueue(c, m).ServeHTTP)
}

func (server *Server) Put(endpoint string, c Handler, m ...Middleware) {
	server.Add(http.MethodPut, endpoint, newQueue(c, m).ServeHTTP)
}

func (server *Server) Patch(endpoint string, c Handler, m ...Middleware) {
	server.Add(http.MethodPatch, endpoint, newQueue(c, m).ServeHTTP)
}

func (server *Server) Delete(endpoint string, c Handler, m ...Middleware) {
	server.Add(http.MethodDelete, endpoint, newQueue(c, m).ServeHTTP)
}

func (server *Server) NewGroup(endpoint string, routeWide ...Middleware) *Group {
	return &Group{endpoint, routeWide, server}
}
