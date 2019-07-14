package middleware

import "net/http"

// A type signature for middlewares
type Middlewares func(http.Handler) http.Handler

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
