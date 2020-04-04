package microframework

import (
	"fmt"
	"net/http"
	"github.com/husobee/vestigo"
)

// purely convience types (less typing)
type Handler func(*ResponseWriter, *Request) error
type Middleware func(Handler) Handler

type Server struct {
	address string
	*vestigo.Router
}

// create a new server object
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
