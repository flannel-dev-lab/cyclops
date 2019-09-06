package router

import (
	"net/http"
)

// methods - a list of methods that are allowed
var (
	methods = map[string]bool{
		http.MethodConnect: true,
		http.MethodDelete:  true,
		http.MethodGet:     true,
		http.MethodHead:    true,
		http.MethodOptions: true,
		http.MethodPatch:   true,
		http.MethodPost:    true,
		http.MethodPut:     true,
		http.MethodTrace:   true,
	}

	// AllowTrace - Globally allow the TRACE method handling within url router. This
	// generally not a good idea to have true in production settings, but excellent for testing.
	AllowTrace = false
)

//validMethod - validate that the http method is valid.
func validMethod(method string) bool {
	_, ok := methods[method]
	return ok
}
