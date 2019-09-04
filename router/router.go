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
	tree        *node
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

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	params := make(map[string]string)

	node, _ := router.tree.searchTree(strings.Split(filepath.Clean(r.URL.Path), "/")[1:], params)
	if handler := node.methods[r.Method]; handler != nil {
		handler(w, r, params)
	} else {
		http.NotFound(w, r)
		return
	}
}
