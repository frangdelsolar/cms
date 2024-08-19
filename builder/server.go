package builder

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// RouteHandler defines a structure for storing route information.
type RouteHandler struct {
	route   string                                   // route is the path for the route. i.e. /users/{id}
	handler func(http.ResponseWriter, *http.Request) // handler is the handler for the route
	name    string                                   // name is the name of the route
}

// Server defines a structure for managing an HTTP server with middleware and routing capabilities.
type Server struct {
	*http.Server                                   // Server is the underlying HTTP server
	middlewares  []func(http.Handler) http.Handler // middlewares is a slice of middleware functions
	routes       []RouteHandler                    // routes is a slice of route handlers
	Router       *mux.Router                       // Router is the router for the server
}

// ServerConfig defines the configuration options for creating a new Server.
type ServerConfig struct {
	Host string // Host is the hostname or IP address to listen on.
	Port string // Port is the port number to listen on.
}

// NewServer creates a new Server instance with the provided configuration.
//
// It checks for missing configuration (Host and Port) and returns an error if necessary.
// Otherwise, it creates a new Gorilla Mux router, sets up the server address and handler,
// and adds a basic logging middleware by default.
func NewServer(config *ServerConfig) (*Server, error) {

	if config == nil {
		return nil, ErrServerConfigNotProvided
	}

	if config.Host == "" {
		config.Host = "localhost"
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	r := mux.NewRouter()

	svr := &Server{
		Server: &http.Server{
			Addr:    config.Host + ":" + config.Port,
			Handler: r,
		},
		middlewares: []func(http.Handler) http.Handler{},
		routes:      []RouteHandler{},
		Router:      r,
	}

	svr.AddMiddleware(loggingMiddleware)

	svr.AddRoute("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Home")
	}, "home")

	return svr, nil
}

// loggingMiddleware is a sample middleware function that logs the request URI.
//
// It takes an http.Handler as input and returns a new http.Handler that wraps the original
// handler and logs the request URI before calling the original handler.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info().Msg(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

// Run starts the server and listens for incoming connections on the configured address.
//
// It logs a message indicating the server is running on the specified port,
// applies all registered middleware to the server's handler,
// and finally calls the underlying http.Server's ListenAndServe method.
func (s *Server) Run() error {
	log.Info().Msgf("Running server on port %s", s.Addr)

	for _, middleware := range s.middlewares {
		s.Handler = middleware(s.Handler)
	}

	for _, route := range s.routes {
		s.Router.HandleFunc(route.route, route.handler).Name(route.name)
	}

	return s.ListenAndServe()
}

// AddMiddleware adds a new middleware function to the server's middleware chain.
//
// Middleware functions are executed sequentially in the order they are added.
// Each middleware function takes an http.Handler as input and returns a new http.Handler
// that can wrap the original handler and perform additional logic before or after
// the original handler is called.
func (s *Server) AddMiddleware(middleware func(http.Handler) http.Handler) {
	s.middlewares = append(s.middlewares, middleware)
}

// AddRoute adds a new route to the server's routing table.
//
// It takes three arguments:
//   - route: The path for the route (e.g., "/", "/users/{id}").
//   - handler: The function to be called when the route is matched.
//   - name: An optional name for the route (useful for generating URLs)
//
// Example:
//
//	AddRoute("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
//	  // Handle user with ID
//	}, "getUser")
//
// url, err := r.Get("getUser").URL("id", "123") =>
// "/users/123"
func (s *Server) AddRoute(route string, handler func(w http.ResponseWriter, r *http.Request), name string) {
	s.routes = append(s.routes, RouteHandler{
		route:   route,
		handler: handler,
		name:    name,
	})
}