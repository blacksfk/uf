/*
A tiny framework implementing routing (via vestigo), middleware support, and error handling.

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
	e := server.Start()

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
	"fmt"
	"github.com/husobee/vestigo"
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
	*vestigo.Router
}

// Server configuration
type Config struct {
	// Address to listen on. Eg. ":6060"
	Address string

	// Logs errors that occur during requests
	ErrorLogger ErrorLogger

	// Logs requests
	AccessLogger AccessLogger

	// CORS configuration as per vestigo.
	// Eg.:
	// &vestigo.CorsAccessControl{
	//   AllowOrigin: []string{"*"},
	//   AllowCredentials: true,
	//   ExposeHeaders: []string{"*"},
	//   MaxAge: time.Hour,
	//   AllowMethods: []string{"GET", "POST", ...},
	//   AllowHeaders: []string{"Content-Type", ...},
	// }
	Cors *vestigo.CorsAccessControl
}

// Create a new server; optionally specifying global middleware.
func NewServer(config *Config, m ...Middleware) *Server {
	router := vestigo.NewRouter()

	if config.Cors != nil {
		router.SetGlobalCors(config.Cors)
	}

	return &Server{config, m, router}
}

// Listen for connections on the previously supplied address
func (s *Server) Start() error {
	fmt.Printf("Starting server on %s\n", s.Config.Address)

	return http.ListenAndServe(s.Config.Address, s)
}

// Bind endpoint to the specified method, append the supplied middleware (if any)
// to the global middleware and create the middleware queue.
func (s *Server) bind(method, endpoint string, c Handler, m []Middleware) {
	m = append(s.GlobalMiddleware, m...)

	s.Add(method, endpoint, newQueue(c, m, s.Config.ErrorLogger, s.Config.AccessLogger).ServeHTTP)
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
