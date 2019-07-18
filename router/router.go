package router

import (
	"net/http"
)

// RegisterRoutes takes in a handler and attaches all the routes to the handler
func RegisterRoutes(handler *http.ServeMux, routes map[string]http.Handler) {
	for k, v := range routes {
		handler.Handle(k, v)
	}
}


func InitializeServer(address string) (*http.ServeMux, *http.Server) {
	handler := http.NewServeMux()

	server := http.Server{
		Addr:    address,
		Handler: handler,
	}
	return handler, &server
}

// InitializeHTTPServer binds a HTTP server on a given address and port
func InitializeHTTPServer(address string) (*http.ServeMux, *http.Server) {
	handler := http.NewServeMux()

	server := http.Server{
		Addr:    address,
		Handler: handler,
	}
	return handler, &server
}
// TODO add TLS support
