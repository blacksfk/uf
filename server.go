/*
A set of utility functions and a router provided by HTTPRouter (https://github.com/julienschmidt/httprouter).

Example:

import (
	uf "github.com/blacksfk/microframework"
	// ...
)

func main() {
	// create a new server
	config := &uf.Config{...}

	// middlewareX and Y are middleware that will be applied to every route defined
	server := uf.NewServer(config, middlewareX, middlewareY, ...)

	// configure CORS or any other settings exported by HTTPRouter
	server.GlobalOPTIONS = func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			h := w.Header()

			h.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			h.Set("Access-Control-Allow-Origin", "example.com")

			// ...
		}

		// reply to options
		w.WriteHeader(http.StatusNoContent)
	}

	// add routes to the server (supports GET, POST, PUT, PATCH, DELETE convenience methods)
	// middleware functions X, Y, A, B, C will be called before the handler in order
	server.Get("/book", handler, middlewareA, middlewareB, middlewareC)

	// ...

	// add route groups to the server
	// middleware functions X, Y, A, B will be called in order before each middleware
	// and handler defined below
	server.Group("/author", middlewareA, middlewareB).

		// middlewareC will be called after X, Y, A, and B but only for GET requests on this route
		Get(author.HandleGet, middlewareC).

		// middlewareD will be called after X, Y, A, and B for the following route definitions
		Middleware(middlewareD).
		Post(author.HandlePost).
		Put(author.HandlePut).

		// middlewareE will be called after X, Y, A, B, and D but only for PATCH
		// requests on this route
		Patch(author.HandlePatch, middlewareE).
		Delete(author.HandleDelete)

	// start the server
	e := http.ListenAndServe("localhost:6060", server)

	// ...
}

func handler(w http.ResponseWriter, r *http.Request) error {
	books := database.ObtainBooks()

	return uf.SendJSON(w, books)
}

func middlewareA(r *http.Request) error {
	// get the auth key and user (somehow)
	key := r.Header.Get("Authorization")
	user := database.FindUser()

	if !user.Valid(key) {
		// user needs to re-authenticate
		return uf.Unauthorized("Invalid login")
	}

	// authenticated
	*r = *r.WithContext(user.ToContext(r.Context()))

	// progress to next handler
	return nil
}
*/
package microframework

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// Extension of http.Handler that returns an error to the framework's error handler
type Handler func(http.ResponseWriter, *http.Request) error

// Middleware functions are called in order after route matching but
// before the Handler.
type Middleware func(*http.Request) error

// Functions implementing this type are supplied an HttpError if
// an error occurs while processing a request.
type ErrorLogger func(error)

// Functions implementing this type are supplied the request and duration
// of the request along with an appropriate unit i.e. m, u, or n.
type AccessLogger func(*http.Request, int64, string)

// Wrapper around vestigo.Router
type Server struct {
	Config           *Config
	GlobalMiddleware []Middleware
	*httprouter.Router
}

// Server configuration
type Config struct {
	// Logs errors that occur during requests
	ErrorLogger ErrorLogger

	// Logs requests
	AccessLogger AccessLogger
}

// Create a new server; optionally specifying global middleware.
func NewServer(config *Config, m ...Middleware) *Server {
	return &Server{config, m, httprouter.New()}
}

// Bind endpoint to the specified method, append the supplied middleware (if any)
// to the global middleware and create the middleware queue.
func (s *Server) bind(method, endpoint string, h Handler, m []Middleware) {
	m = append(s.GlobalMiddleware, m...)

	s.Handler(method, endpoint, newQueue(h, m, s.Config))
}

// Bind endpoint to support GET requests.
func (s *Server) Get(endpoint string, c Handler, m ...Middleware) {
	s.bind(http.MethodGet, endpoint, c, m)
}

// Bind endpoint to support POST requests.
func (s *Server) Post(endpoint string, c Handler, m ...Middleware) {
	s.bind(http.MethodPost, endpoint, c, m)
}

// Bind endpoint to support PUT requests.
func (s *Server) Put(endpoint string, c Handler, m ...Middleware) {
	s.bind(http.MethodPut, endpoint, c, m)
}

// Bind endpoint to support PATCH requests.
func (s *Server) Patch(endpoint string, c Handler, m ...Middleware) {
	s.bind(http.MethodPatch, endpoint, c, m)
}

// Bind endpoint to support DELETE requests.
func (s *Server) Delete(endpoint string, c Handler, m ...Middleware) {
	s.bind(http.MethodDelete, endpoint, c, m)
}

// Append (or set if not existing) middleware to apply to all routes.
func (s *Server) AddGlobalMiddleware(m ...Middleware) {
	s.GlobalMiddleware = append(s.GlobalMiddleware, m...)
}

// Create a group to bind multiple HTTP verbs to a single endpoint concisely
func (s *Server) NewGroup(endpoint string, routeWide ...Middleware) *Group {
	return &Group{endpoint, routeWide, s}
}
