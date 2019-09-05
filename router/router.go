package router

import (
	"net/http"
	"path/filepath"
	"strings"
)

// Handle type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, Handle(f) is a
// Handler that calls f.
type Handle func(http.ResponseWriter, *http.Request, map[string]string)

type Router struct {
	tree          *node
	staticPath    string
	staticHandler http.Handler
}

func New() *Router {
	node := node{component: "/", isNamedParam: false, methods: make(map[string]Handle)}
	return &Router{tree: &node}
}

func (router *Router) Get(path string, handle Handle) {
	router.tree.addNode(http.MethodGet, path, handle)
}

func (router *Router) Post(path string, handle Handle) {
	router.tree.addNode(http.MethodPost, path, handle)
}

func (router *Router) Put(path string, handle Handle) {
	router.tree.addNode(http.MethodPut, path, handle)
}

func (router *Router) Patch(path string, handle Handle) {
	router.tree.addNode(http.MethodPatch, path, handle)
}

func (router *Router) Delete(path string, handle Handle) {
	router.tree.addNode(http.MethodDelete, path, handle)
}

func (router *Router) Head(path string, handle Handle) {
	router.tree.addNode(http.MethodHead, path, handle)
}

func (router *Router) Trace(path string, handle Handle) {
	router.tree.addNode(http.MethodTrace, path, handle)
}

func (router *Router) Options(path string, handle Handle) {
	router.tree.addNode(http.MethodOptions, path, handle)
}

func (router *Router) Connect(path string, handle Handle) {
	router.tree.addNode(http.MethodConnect, path, handle)
}

func (router *Router) Custom(method, path string, handle Handle) {
	router.tree.addNode(method, path, handle)
}

func (router *Router) RegisterStatic(directoryPath, servePath string) {
	fs := http.FileServer(FileSystem{http.Dir(directoryPath)})
	router.staticHandler = http.StripPrefix(servePath, fs)
	router.staticPath = servePath
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := make(map[string]string)

	if strings.Contains(r.URL.Path, router.staticPath) && router.staticPath != "" {
		router.staticHandler.ServeHTTP(w, r)
	} else {
		node, _ := router.tree.searchTree(strings.Split(filepath.Clean(r.URL.Path), "/")[1:], params)
		if handler := node.methods[r.Method]; handler != nil {
			handler(w, r, params)
		} else {
			http.NotFound(w, r)
			return
		}
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
