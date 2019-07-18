package middleware

import (
	"net/http"
	"strings"
)

// DefaultHeaders lets you manage a set of default headers as per mozilla spec
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Expose-Headers
type DefaultHeaders struct {
	// CacheControl general-header field is used to specify directives for caching mechanisms in both requests
	// and responses.
	CacheControl string
	// ContentLanguage entity header is used to describe the language(s) intended for the audience, so that it allows a user to
	// differentiate according to the users' own preferred language. Default is en-US
	ContentLanguage []string
	// ContentType entity header is used to indicate the media type of the resource.
	ContentType string
	// Expires header contains the date/time after which the response is considered stale.
	Expires string
	// LastModified response HTTP header contains the date and time at which the origin server believes the resource was last
	// modified. It is used as a validator to determine if a resource received or stored is the same.
	// Format: <day-name>, <day> <month> <year> <hour>:<minute>:<second> GMT
	LastModified string
}

// SetDefaultHeaders will set certain default headers specified by the user
func (defaultHeaders *DefaultHeaders) SetDefaultHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		if defaultHeaders.CacheControl != "" {
			w.Header().Set("Cache-Control", defaultHeaders.CacheControl)
		}

		if len(defaultHeaders.ContentLanguage) > 0 {
			w.Header().Set("Content-Language", strings.Join(defaultHeaders.ContentLanguage, ", "))
		}

		if defaultHeaders.ContentType != "" {
			w.Header().Set("Content-Type", defaultHeaders.ContentType)
		}

		if defaultHeaders.Expires != "" {
			w.Header().Set("Expires", defaultHeaders.Expires)
		}

		if defaultHeaders.LastModified != "" {
			w.Header().Set("Last-Modified", defaultHeaders.LastModified)
		}
		h.ServeHTTP(w, request)

	})
}

// SetSecureHeaders sets some default security headers
func SetSecureHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		h.ServeHTTP(w, request)
	})
}
