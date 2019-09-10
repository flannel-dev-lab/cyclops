package router

import (
	"github.com/flannel-dev-lab/cyclops/input"
	"net/http"
	"strings"
)

const (
	stype ntype = iota
	ptype
	mtype
)

type (
	ntype    uint8
	children []*node
)

// Router - The main vestigo router data structure
type Router struct {
	root          *node
	staticHandler http.Handler
	staticPath    string
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

// NewRouter - Create a new vestigo router
func New() *Router {
	return &Router{
		root: &node{
			resource: newResource(),
		},
	}
}

// GetMatchedPathTemplate - get the path template from the url in the request
func (r *Router) GetMatchedPathTemplate(req *http.Request) string {
	p, _ := r.find(req)
	return p
}

// ServeHTTP - implementation of a http.Handler, making Router a http.Handler
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if strings.Contains(req.URL.Path, r.staticPath) && r.staticPath != "" {
		r.staticHandler.ServeHTTP(w, req)
	} else {
		h := r.Find(req)
		h(w, req)
	}

}

// Get - Helper method to add HTTP GET Method to router
func (r *Router) Get(path string, handler http.HandlerFunc) {
	r.Add(http.MethodGet, path, handler)
}

// Post - Helper method to add HTTP POST Method to router
func (r *Router) Post(path string, handler http.HandlerFunc) {
	r.Add(http.MethodPost, path, handler)
}

// Connect - Helper method to add HTTP CONNECT Method to router
func (r *Router) Connect(path string, handler http.HandlerFunc) {
	r.Add(http.MethodConnect, path, handler)
}

// Delete - Helper method to add HTTP DELETE Method to router
func (r *Router) Delete(path string, handler http.HandlerFunc) {
	r.Add(http.MethodDelete, path, handler)
}

// Patch - Helper method to add HTTP PATCH Method to router
func (r *Router) Patch(path string, handler http.HandlerFunc) {
	r.Add(http.MethodPatch, path, handler)
}

// Put - Helper method to add HTTP PUT Method to router
func (r *Router) Put(path string, handler http.HandlerFunc) {
	r.Add(http.MethodPut, path, handler)
}

// Trace - Helper method to add HTTP TRACE Method to router
func (r *Router) Trace(path string, handler http.HandlerFunc) {
	r.Add(http.MethodTrace, path, handler)
}

// Head - Helper method to add HTTP HEAD Method to router
func (r *Router) Head(path string, handler http.HandlerFunc) {
	r.Add(http.MethodHead, path, handler)
}

// Add - Add a method/handler combination to the router
func (r *Router) Add(method, path string, h http.HandlerFunc) {
	r.add(method, path, h)
}

// Add - Add a method/handler combination to the router
func (r *Router) add(method, path string, h http.HandlerFunc) {
	pnames := make(pNames)
	pnames[method] = []string{}
	for i, l := 0, len(path); i < l; i++ {
		if path[i] == ':' {
			j := i + 1

			r.insert(method, path[:i], nil, stype, nil)
			for ; i < l && path[i] != '/'; i++ {
			}

			pnames[method] = append(pnames[method], path[j:i])
			path = path[:j] + path[i:]
			i, l = j, len(path)

			if i == l {
				r.insert(method, path[:i], h, ptype, pnames)
				return
			}
			r.insert(method, path[:i], nil, ptype, pnames)
		} else if path[i] == '*' {
			r.insert(method, path[:i], nil, stype, nil)
			pnames[method] = append(pnames[method], "_name")
			r.insert(method, path[:i+1], h, mtype, pnames)
			return
		}
	}

	r.insert(method, path, h, stype, pnames)
}

// Find - Find A route within the router tree
func (r *Router) Find(req *http.Request) (h http.HandlerFunc) {
	_, h = r.find(req)
	return
}

func (r *Router) find(req *http.Request) (prefix string, h http.HandlerFunc) {
	// get tree base node from the router
	cn := r.root

	h = notFoundHandler

	if !validMethod(req.Method) {
		// if the method is completely invalid
		h = methodNotAllowedHandler(cn.resource.allowedMethods)
		return
	}

	var (
		search          = req.URL.Path
		c               *node // Child node
		n               int   // Param counter
		collectedPnames []string
	)

	// Search order static > param > match-any
	for {

		if search == "" {
			if cn.resource != nil {
				// Found route, check if method is applicable
				handler, allowedMethods := cn.resource.GetMethodHandler(req.Method)
				if handler == nil {
					if allowedMethods != "" {
						// route is valid, but method is not allowed, 405
						h = methodNotAllowedHandler(allowedMethods)
					}
					return
				}
				h = handler
				for i, v := range collectedPnames {
					if len(cn.pnames[req.Method]) > i {
						input.AddParam(req, cn.pnames[req.Method][i], v)
					}
				}

				brokenPrefix := strings.Split(prefix, "/")
				prefix = ""
				k := 0
				for _, v := range brokenPrefix {
					if v != "" {
						prefix += "/"
						if v == ":" {
							if pnames, ok := cn.pnames[req.Method]; ok {
								prefix += v + pnames[k]
							}
							k++
						} else {
							prefix += v
						}
					}
				}
			}
			return
		}

		pl := 0 // Prefix length
		l := 0  // LCP length

		if cn.label != ':' {
			sl := len(search)
			pl = len(cn.prefix)
			prefix += cn.prefix

			// LCP
			max := pl
			if sl < max {
				max = sl
			}
			for ; l < max && search[l] == cn.prefix[l]; l++ {
			}
		}

		if l == pl {
			// Continue search
			search = search[l:]

			if search == "" && cn != nil && cn.parent != nil && cn.resource.allowedMethods == "" {
				parent := cn.parent
				search = cn.prefix
				for parent != nil {
					if sib := parent.findChildWithLabel('*'); sib != nil {
						search = parent.prefix + search
						cn = parent
						goto MatchAny
					}
					parent = parent.parent
				}
			}

		}

		if search == "" {
			// TODO: Needs improvement
			if cn.findChildWithType(mtype) == nil {
				continue
			}
			// Empty value
			goto MatchAny
		}

		// Static node
		c = cn.findChild(search, stype)
		if c != nil {
			cn = c
			continue
		}
		// Param node
	Param:

		c = cn.findChildWithType(ptype)
		if c != nil {
			cn = c

			i, l := 0, len(search)
			for ; i < l && search[i] != '/'; i++ {
			}

			collectedPnames = append(collectedPnames, search[0:i])
			prefix += ":"
			n++
			search = search[i:]

			if len(cn.children) == 0 && len(search) != 0 {
				return
			}

			continue
		}

		// Match-any node
	MatchAny:
		//		c = cn.getChild()
		c = cn.findChildWithType(mtype)
		if c != nil {
			cn = c
			collectedPnames = append(collectedPnames, search)
			search = "" // End search
			continue
		}
		// last ditch effort to match on wildcard (issue #8)
		var tmpsearch = search
		for {
			if cn != nil && cn.parent != nil && cn.prefix != ":" {
				tmpsearch = cn.prefix + tmpsearch
				cn = cn.parent
				if cn.prefix == "/" {
					var sib = cn.findChildWithLabel(':')
					if sib != nil {
						search = tmpsearch
						goto Param
					}
					if sib := cn.findChildWithLabel('*'); sib != nil {
						search = tmpsearch
						goto MatchAny
					}
				}
			} else {
				break
			}
		}

		// Not found
		return
	}
}

// insert - insert a route into the router tree
func (r *Router) insert(method, path string, h http.HandlerFunc, t ntype, pnames pNames) {
	// Adjust max param

	cn := r.root

	if !validMethod(method) && method != "CORS" {
		panic("invalid method")
	}
	search := path

	for {
		sl := len(search)
		pl := len(cn.prefix)
		l := 0

		// LCP
		max := pl
		if sl < max {
			max = sl
		}
		for ; l < max && search[l] == cn.prefix[l]; l++ {
		}

		if cn.pnames == nil {
			cn.pnames = make(pNames)
		}

		if l == 0 {
			// At root node
			cn.label = search[0]
			cn.prefix = search
			if h != nil {
				cn.typ = t
				cn.resource = newResource()
				if method != "CORS" {
					cn.resource.AddMethodHandler(method, h)
				}
				if method == "GET" {
					cn.pnames["HEAD"] = pnames[method]
				}
				cn.pnames[method] = pnames[method]
			}
		} else if l < pl {
			// Split node
			nr := newResource()
			cn.resource.CopyTo(nr)

			n := newNode(cn.typ, cn.prefix[l:], cn, cn.children, nr, cn.pnames)
			for i := 0; i < len(n.children); i++ {
				n.children[i].parent = n
			}

			// Reset parent node
			cn.typ = stype
			cn.label = cn.prefix[0]
			cn.prefix = cn.prefix[:l]
			cn.children = nil
			cn.resource = newResource()
			cn.pnames = make(pNames)

			cn.addChild(n)

			if l == sl {
				// At parent node
				cn.typ = t

				if method != "CORS" {
					cn.resource.AddMethodHandler(method, h)
				}
				if method == "GET" {
					cn.pnames["HEAD"] = pnames[method]
				}
				cn.pnames[method] = pnames[method]
			} else {
				// Create child node
				nr := newResource()
				if method != "CORS" {
					nr.AddMethodHandler(method, h)
				}
				cn.pnames[method] = pnames[method]
				n = newNode(t, search[l:], cn, nil, nr, cn.pnames)
				cn.addChild(n)
			}
		} else if l < sl {
			search = search[l:]
			c := cn.findChildWithLabel(search[0])
			if c != nil {
				// Go deeper
				cn = c
				continue
			}
			// Create child node
			nr := newResource()
			if method != "CORS" {
				nr.AddMethodHandler(method, h)
			}
			n := newNode(t, search, cn, nil, nr, pnames)
			cn.addChild(n)

			cn.resource.Clean()
			n.resource.Clean()

		} else {
			// Node already exists
			if h != nil {
				// add the handler to the node's map of methods to handlers

				if method != "CORS" {
					cn.resource.AddMethodHandler(method, h)
				}
				if method == "GET" {
					cn.pnames["HEAD"] = pnames[method]
				}
				cn.pnames[method] = pnames[method]
			}
		}
		return
	}
}

func (r *Router) RegisterStatic(directoryPath, servePath string) {
	fs := http.FileServer(FileSystem{http.Dir(directoryPath)})
	r.staticHandler = http.StripPrefix(servePath, fs)
	r.staticPath = servePath
}
