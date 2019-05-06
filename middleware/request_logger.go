package middleware

import (
	"log"
	"net/http"
)

// Middleware to log access logs
func RequestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		defer log.Printf("%s %s %s %s %s",
			request.RemoteAddr,
			request.Method,
			request.URL.Path,
			request.Host,
			request.Proto)
		h.ServeHTTP(w, request)
	})
}
