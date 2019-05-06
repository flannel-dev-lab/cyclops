package middleware

import "net/http"

// Validates the content type, If doesn't match it will return 415 else serves the handler
func RequestContentTypeFilter(h http.Handler, contentType string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Content-Type") != contentType {
			http.Error(w, "unsupported content type", http.StatusUnsupportedMediaType)
			return // means request will not go to original handler
		}
		h.ServeHTTP(w, request)
	})
}
