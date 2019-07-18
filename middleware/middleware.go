// middleware package includes different middleware available by default with cyclops, cyclops is made in such a way
// that  it is easy for developers to plug custom middleware as well, the only thing that the developer need to do is
// write a middleware that takes in a http.Handler and returns a http.Handler, once the middleware is complete pass it
// to NewChain method to start using it
package middleware

import "net/http"

// A type signature for middleware
type Middlewares func(http.Handler) http.Handler

// Chain contains a slice of middleware for the request
type Chain struct {
	middlewareHandlers []Middlewares
}

func NewChain(middlewares ...Middlewares) *Chain {
	return &Chain{middlewareHandlers: middlewares}
}

func (chain *Chain) Then(handler http.Handler) http.Handler {
	for _, middleware := range chain.middlewareHandlers {
		handler = middleware(handler)
	}
	handler = SetHeaders(RequestLogger(PanicHandler(handler)))
	return handler
}
