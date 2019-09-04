package cyclops

import (
	"net/http"
	"strings"
)

// RegisterRoutes takes in a handler and attaches all the routes to the handler
func RegisterRoutes(handler *http.ServeMux, routes map[string]http.Handler) {
	for k, v := range routes {
		handler.Handle(k, v)
	}
}

// FileSystem custom file system handler
type FileSystem struct {
	fs http.FileSystem
}

// Open opens file
func (fs FileSystem) Open(path string) (http.File, error) {
	f, err := fs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := fs.fs.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}

// RegisterStatic helps in handling static assets. directoryPath takes in the path to your static assets, servePath
// will take in a path to handler on which the files should be served.
func RegisterStatic(handler *http.ServeMux, directoryPath, servePath string) {
	fs := http.FileServer(FileSystem{http.Dir(directoryPath)})

	handler.Handle(servePath, http.StripPrefix(strings.TrimRight(servePath, "/"), fs))
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
