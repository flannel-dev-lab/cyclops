package middleware

import "net/http"

// Sets the response headers to a response. By default, the following headers are loaded
// "X-Content-Type-Options", "nosniff"
// "X-Frame-Options", "deny"
// "Content-Type", "application/json"
// "X-XSS-Protection", "1; mode=block"
func SetHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		h.ServeHTTP(w, request)
	})
}
