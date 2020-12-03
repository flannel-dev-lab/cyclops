// middleware package includes different middleware available by default with cyclops, cyclops is made in such a way
// that  it is easy for developers to plug custom middleware as well, the only thing that the developer need to do is
// write a middleware that takes in a http.Handler and returns a http.Handler, once the middleware is complete pass it
// to NewChain method to start using it
package middleware

import (
	"net/http"
)

// A type signature for middleware
//type Middlewares func(http.Handler) http.Handler
type Middlewares func(http.HandlerFunc) http.HandlerFunc

// Chain contains a slice of middleware for the request
type Chain struct {
	middlewareHandlers []Middlewares
}

// NewChain takes a variable number of middleware's and adds them to chain and returns a pointer to Chain
func NewChain(middlewares ...Middlewares) *Chain {
	return &Chain{middlewareHandlers: middlewares}
}

// Then will take in your handler that need to be executed with the requested path and chains all the middleware's that
// are  specified in Chain, the note here is that middleware's are chained in the order they are specified,
// so take care of adding middleware's in appropriate order
func (chain *Chain) Then(handler http.HandlerFunc) http.HandlerFunc {
	for _, middleware := range chain.middlewareHandlers {
		handler = middleware(handler)
	}
	handler = AccessLogger(PanicHandler(handler))
	return handler
}
