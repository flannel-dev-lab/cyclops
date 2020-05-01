package router

import (
	"io"
	"net/http"
)

var (
	// traceHandler - Generic Trace Handler to echo back input
	traceHandler = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "message/http")
		w.WriteHeader(http.StatusOK)
		if r.Body == nil {
			_, _ = w.Write([]byte{})
			return
		}
		defer r.Body.Close()
		io.Copy(w, r.Body)
	}
	// headHandler - Generic Head Handler to return header information
	headHandler = func(f http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			fakeWriter := &headResponseWriter{}
			// issue 23 - nodes that do not have handlers should not be called when HEAD
			// is called
			if f != nil {
				f(fakeWriter, r)
				for k, v := range fakeWriter.Header() {
					for _, vv := range v {
						w.Header().Add(k, vv)
					}
				}
				w.WriteHeader(fakeWriter.Code)
				w.Write([]byte(""))
			} else {
				notFoundHandler(w, r)
			}
		}
	}

	// methodNotAllowedHandler - Generic Handler to handle when method isn't allowed for a resource
	methodNotAllowedHandler = func(allowedMethods string) func(w http.ResponseWriter, r *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Allow", allowedMethods)
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		}
	}
	// notFoundHandler - Generic Handler to handle when resource isn't found
	notFoundHandler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
	}
)

// headResponseWriter - implementation of http.ResponseWriter for headHandler
type headResponseWriter struct {
	HeaderMap   http.Header
	Code        int
	wroteHeader bool
}

func (hrw *headResponseWriter) Header() http.Header {
	if hrw.HeaderMap == nil {
		hrw.HeaderMap = make(http.Header)
	}
	return hrw.HeaderMap
}

func (hrw *headResponseWriter) Write([]byte) (int, error) {
	// Mirror http.ResponseWriter: "If WriteHeader has not yet been called,
	// Write calls WriteHeader(http.StatusOK) before writing the data."
	if !hrw.wroteHeader {
		hrw.WriteHeader(http.StatusOK)
	}

	return 0, nil
}

func (hrw *headResponseWriter) WriteHeader(status int) {
	hrw.wroteHeader = true
	hrw.Code = status
}

// resource - internal structure for specifying which handlers belong to a particular route
type resource struct {
	Connect        http.HandlerFunc
	Delete         http.HandlerFunc
	Get            http.HandlerFunc
	Patch          http.HandlerFunc
	Post           http.HandlerFunc
	Put            http.HandlerFunc
	Trace          http.HandlerFunc
	Head           http.HandlerFunc
	allowedMethods string
}

// newResource - create a new resource, and give it sane default values
func newResource() *resource {
	return &resource{
		allowedMethods: "",
	}
}

// CopyTo - Copy the Resource to another Resource passed in by reference
func (h *resource) CopyTo(v *resource) {
	v.Get = h.Get
	v.Connect = h.Connect
	v.Delete = h.Delete
	v.Get = h.Get
	v.Patch = h.Patch
	v.Post = h.Post
	v.Put = h.Put
	v.Trace = h.Trace
	v.Head = h.Head
	v.allowedMethods = h.allowedMethods
}

// addToAllowedMethods - Add a method to the allowed methods for this route
func (h *resource) addToAllowedMethods(method string) {
	if h.allowedMethods == "" {
		h.allowedMethods = method
	} else {
		h.allowedMethods = h.allowedMethods + ", " + method
	}
}

// Clean - Clean up allowed methods based on funcs
func (h *resource) Clean() {
	h.allowedMethods = ""
	hasOneMethod := false
	if h.Get != nil {
		h.addToAllowedMethods(http.MethodGet)
		h.addToAllowedMethods(http.MethodHead)
		h.Head = headHandler(h.Get)
		hasOneMethod = true
	}
	if h.Put != nil {
		h.addToAllowedMethods(http.MethodPut)
		hasOneMethod = true
	}
	if h.Post != nil {
		h.addToAllowedMethods(http.MethodPost)
		hasOneMethod = true
	}
	if h.Patch != nil {
		h.addToAllowedMethods(http.MethodPatch)
		hasOneMethod = true
	}
	if h.Delete != nil {
		h.addToAllowedMethods(http.MethodDelete)
		hasOneMethod = true
	}
	if h.Connect != nil {
		h.addToAllowedMethods(http.MethodConnect)
		hasOneMethod = true
	}
	if h.Head != nil {
		h.addToAllowedMethods(http.MethodHead)
		hasOneMethod = true
	}
	if hasOneMethod && AllowTrace {
		h.addToAllowedMethods(http.MethodTrace)
		h.Trace = traceHandler
	}
}

// AddMethodHandler - Add a method/handler pair to the resource structure
func (h *resource) AddMethodHandler(method string, handler http.HandlerFunc) {
	l := len(method)
	firstChar := method[0]
	secondChar := method[1]
	if h != nil {
		if AllowTrace {
			h.addToAllowedMethods(http.MethodTrace)
			h.Trace = traceHandler
		}
		if l == 3 {
			if uint16(firstChar)<<8|uint16(secondChar) == 0x4745 {
				h.addToAllowedMethods(method)
				h.addToAllowedMethods(http.MethodHead)
				h.Get = handler
				h.Head = headHandler(handler)
			}
			if uint16(firstChar)<<8|uint16(secondChar) == 0x5055 {
				h.addToAllowedMethods(method)
				h.Put = handler
			}
		} else if l == 4 {
			if uint16(firstChar)<<8|uint16(secondChar) == 0x504f {
				h.addToAllowedMethods(method)
				h.Post = handler
			}
			if uint16(firstChar)<<8|uint16(secondChar) == 0x4845 {
				h.addToAllowedMethods(method)
				h.Head = handler
			}
		} else if l == 5 {
			if uint16(firstChar)<<8|uint16(secondChar) == 0x5452 {
				h.addToAllowedMethods(method)
				h.Trace = handler
			}
			if uint16(firstChar)<<8|uint16(secondChar) == 0x5041 {
				h.addToAllowedMethods(method)
				h.Patch = handler
			}
		} else if l >= 6 {
			if uint16(firstChar)<<8|uint16(secondChar) == 0x4445 {
				h.addToAllowedMethods(method)
				h.Delete = handler
			}
			if uint16(firstChar)<<8|uint16(secondChar) == 0x434f {
				h.addToAllowedMethods(method)
				h.Connect = handler
			}
		}
	}
}

// GetMethodHandler - Get a method/handler pair from the resource structure
func (h *resource) GetMethodHandler(method string) (http.HandlerFunc, string) {
	l := len(method)
	firstChar := method[0]
	secondChar := method[1]
	if l == 3 {
		if uint16(firstChar)<<8|uint16(secondChar) == 0x4745 {
			return h.Get, h.allowedMethods
		}
		if uint16(firstChar)<<8|uint16(secondChar) == 0x5055 {
			return h.Put, h.allowedMethods
		}
	} else if l == 4 {
		if uint16(firstChar)<<8|uint16(secondChar) == 0x504f {
			return h.Post, h.allowedMethods
		}
		if uint16(firstChar)<<8|uint16(secondChar) == 0x4845 {
			return h.Head, h.allowedMethods
		}
	} else if l == 5 {
		if uint16(firstChar)<<8|uint16(secondChar) == 0x5452 {
			return h.Trace, h.allowedMethods
		}
		if uint16(firstChar)<<8|uint16(secondChar) == 0x5041 {
			return h.Patch, h.allowedMethods
		}
	} else if l >= 6 {
		if uint16(firstChar)<<8|uint16(secondChar) == 0x4445 {
			return h.Delete, h.allowedMethods
		}
		if uint16(firstChar)<<8|uint16(secondChar) == 0x434f {
			return h.Connect, h.allowedMethods
		}
	}
	return nil, h.allowedMethods
}
