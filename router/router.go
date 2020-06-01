package router

import (
	"log"
	"net/http"
	"strings"
)

var (
	CyclopsNotFoundHandler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}

	CyclopsMethodNotAllowedHandler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
)

// Router is a struct which handles dispatching requests to different handlers
type Router struct {
	tree          *node
	staticHandler http.Handler
	staticPath    string
	// StripTrailingSlashOnRequest removes a trailing slash when registering routes
	StripTrailingSlashOnRegisteringHandlers bool
	// NotFoundHandler allows you to pass in a custom NotFoundHandler when a handler
	// is not found
	NotFoundHandler http.HandlerFunc
	// MethodNotAllowedHandler allows you to pass in a custom MethodNotAllowedHandler
	// when a method is not allowed
	MethodNotAllowedHandler http.HandlerFunc
}

func New(
	stripTrailingSlashOnRegisteringHandlers bool,
	notFoundHandler,
	methodNotAllowedHandler http.HandlerFunc) *Router {

	node := &node{
		children:     []*node{},
		component:    "/",
		isNamedParam: false,
		methods:      make(map[string]http.HandlerFunc),
	}

	router := &Router{
		tree:                                    node,
		StripTrailingSlashOnRegisteringHandlers: stripTrailingSlashOnRegisteringHandlers,
	}

	if notFoundHandler == nil {
		router.NotFoundHandler = CyclopsNotFoundHandler
	} else {
		router.NotFoundHandler = notFoundHandler
	}

	if methodNotAllowedHandler == nil {
		router.MethodNotAllowedHandler = CyclopsMethodNotAllowedHandler
	} else {
		router.MethodNotAllowedHandler = methodNotAllowedHandler
	}

	return router
}

// Get - Helper method to add HTTP GET Method to router
func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.add(http.MethodGet, path, handler)
}

// Post - Helper method to add HTTP POST Method to router
func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.add(http.MethodPost, path, handler)
}

// Connect - Helper method to add HTTP CONNECT Method to router
func (r *Router) Connect(path string, handler http.HandlerFunc) {
	r.add(http.MethodConnect, path, handler)
}

// Delete - Helper method to add HTTP DELETE Method to router
func (r *Router) Delete(path string, handler http.HandlerFunc) {
	r.add(http.MethodDelete, path, handler)
}

// Patch - Helper method to add HTTP PATCH Method to router
func (r *Router) Patch(path string, handler http.HandlerFunc) {
	r.add(http.MethodPatch, path, handler)
}

// Put - Helper method to add HTTP PUT Method to router
func (r *Router) Put(path string, handler http.HandlerFunc) {
	r.add(http.MethodPut, path, handler)
}

// Trace - Helper method to add HTTP TRACE Method to router
func (r *Router) Trace(path string, handler http.HandlerFunc) {
	r.add(http.MethodTrace, path, handler)
}

// Head - Helper method to add HTTP HEAD Method to router
func (r *Router) Head(path string, handler http.HandlerFunc) {
	r.add(http.MethodHead, path, handler)
}

// Options - Helper method to add HTTP OPTIONS Method to router
func (r *Router) Options(path string, handler http.HandlerFunc) {
	r.add(http.MethodOptions, path, handler)
}

func (r *Router) add(method, path string, handler http.HandlerFunc) {
	if method == "" {
		panic("method must not be empty")
	}

	if len(path) < 1 || path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if handler == nil {
		panic("handler must not be nil")
	}

	if r.StripTrailingSlashOnRegisteringHandlers && path != "/" {
		path = strings.TrimSuffix(path, "/")
	}

	r.tree.addNode(method, path, handler)
}

func (r *Router) find(req *http.Request) (http.HandlerFunc, error) {
	_ = req.ParseForm()

	params := req.Form

	node, _ := r.tree.traverse(strings.Split(req.URL.Path, "/")[1:], params)
	if handler := node.methods[req.Method]; handler != nil {
		q := req.URL.Query()

		for key, values := range params {
			temp := values
			q.Del(key)

			for _, value := range temp {
				q.Add(key, value)
			}
		}

		req.URL.RawQuery = q.Encode()

		return handler, nil
	} else {
		if len(node.methods) == 0 {
			return r.NotFoundHandler, nil
		} else {
			return r.MethodNotAllowedHandler, nil
		}
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if strings.Contains(req.URL.Path, r.staticPath) && r.staticPath != "" {
		r.staticHandler.ServeHTTP(w, req)
	} else {
		handler, _ := r.find(req)
		handler(w, req)
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
		log.Println(err)
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := fs.fs.Open(index); err != nil {
			log.Println(err)
			return nil, err
		}
	}

	return f, nil
}

// RegisterStatic registers a static directory to serve on servePath
func (r *Router) RegisterStatic(directoryPath, servePath string) {
	fs := http.FileServer(FileSystem{http.Dir(directoryPath)})
	r.staticHandler = http.StripPrefix(servePath, fs)
	r.staticPath = servePath
}
