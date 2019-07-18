package middleware

// Contains all the CORS configurations
type CORS struct {
	// AllowedOrigin Specifies which origin should be allowed, if you want to allow all use *
	AllowedOrigin string

	// AllowedCredentials indicates whether the response to the request can be exposed when the credentials flag is true.
	// The only valid value for this header is true (case-sensitive). If you don't need credentials,
	// omit this header entirely (rather than setting its value to false).
	AllowedCredentials bool

	// TODO Convert this to space seperated strings
	// AllowedHeaders is used in response to a preflight request which includes the Access-Control-Request-Headers
	// to indicate which HTTP headers can be used during the actual request.
	AllowedHeaders []string

	// TODO Convert this to space seperated strings
	// AllowedMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (HEAD, GET and POST).
	AllowedMethods []string

	// TODO Convert this to space seperated strings
	// ExposedHeaders indicates which headers can be exposed as part of the response by listing their names.
	ExposedHeaders []string

	// MaxAge indicates how long the results of a preflight request (that is the information contained in the
	// Access-Control-Allow-Methods and Access-Control-Allow-Headers headers) can be cached.
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age
	MaxAge int
}

func (cors *CORS) New() {
	if cors.AllowedOrigin == "" {
		cors.AllowedOrigin = "*"
	}
}