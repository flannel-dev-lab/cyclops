package middleware

import (
	"github.com/flannel-dev-lab/cyclops/router"
	"net/http"
	"strconv"
	"strings"
)

// Contains all the CORS configurations
type CORS struct {
	// AllowedOrigin Specifies which origin should be allowed, if you want to allow all use *
	AllowedOrigin string

	// AllowedCredentials indicates whether the response to the request can be exposed when the credentials flag is true.
	// The only valid value for this header is true (case-sensitive). If you don't need credentials,
	// omit this header entirely (rather than setting its value to false).
	AllowedCredentials bool

	// AllowedHeaders is used in response to a pre-flight request which includes the Access-Control-Request-Headers
	// to indicate which HTTP headers can be used during the actual request.
	AllowedHeaders []string

	// AllowedMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (HEAD, GET and POST).
	AllowedMethods []string

	// ExposedHeaders indicates which headers can be exposed as part of the response by listing their names.
	ExposedHeaders []string

	// MaxAge indicates how long the results of a pre-flight request (that is the information contained in the
	// Access-Control-Allow-Methods and Access-Control-Allow-Headers headers) can be cached.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age
	MaxAge int
}

// CORSHandler handles the simple and pre-flight requests
func (cors CORS) CORSHandler(h router.Handle) router.Handle {
	return func(w http.ResponseWriter, request *http.Request, params map[string]string) {
		if request.Method == http.MethodOptions {
			// Handling pre-flight requests
			cors.handlePreflight(w, request)
			h(w, request, params)

		} else {
			// Handling simple request
			cors.handleSimple(w, request)
			h(w, request, params)
		}
	}
}

// handleSimple handles the simple requests
func (cors CORS) handleSimple(w http.ResponseWriter, request *http.Request) {
	if cors.AllowedOrigin != "" {
		w.Header().Set("Access-Control-Allow-Origin", cors.AllowedOrigin)
	}

	if cors.AllowedCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	if len(cors.AllowedHeaders) > 0 {
		for idx, header := range cors.AllowedHeaders {
			cors.AllowedHeaders[idx] = http.CanonicalHeaderKey(header)
		}
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(cors.AllowedHeaders, ", "))
	}

	if len(cors.AllowedMethods) > 0 {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(cors.AllowedMethods, ", "))
	}

	if len(cors.ExposedHeaders) > 0 {
		for idx, header := range cors.ExposedHeaders {
			cors.ExposedHeaders[idx] = http.CanonicalHeaderKey(header)
		}
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(cors.ExposedHeaders, ", "))
	}

	w.Header().Set("Access-Control-Max-Age", strconv.Itoa(cors.MaxAge))
}

// handlePreflight handles the pre-flight requests
func (cors CORS) handlePreflight(w http.ResponseWriter, request *http.Request) {
	if cors.AllowedOrigin != "" {
		w.Header().Set("Access-Control-Allow-Origin", cors.AllowedOrigin)
	}

	if cors.AllowedCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	if len(cors.AllowedHeaders) > 0 {
		for idx, header := range cors.AllowedHeaders {
			cors.AllowedHeaders[idx] = http.CanonicalHeaderKey(header)
		}
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(cors.AllowedHeaders, ", "))
	}

	if len(cors.AllowedMethods) > 0 {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(cors.AllowedMethods, ", "))
	}

	if len(cors.ExposedHeaders) > 0 {
		for idx, header := range cors.ExposedHeaders {
			cors.ExposedHeaders[idx] = http.CanonicalHeaderKey(header)
		}
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(cors.ExposedHeaders, ", "))
	}

	w.Header().Set("Access-Control-Max-Age", strconv.Itoa(cors.MaxAge))
}
