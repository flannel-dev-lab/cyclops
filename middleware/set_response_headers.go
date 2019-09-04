package middleware

import (
	"github.com/flannel-dev-lab/cyclops/router"
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

type SecureHeaders struct {
	// Sets the X-XSS-Protection Header. Valid values are https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-XSS-Protection
	// Default is 1; mode=block
	XSSProtection string
	// ContentTypeOptions provides protection against overriding Content-Type
	// header by setting the `X-Content-Type-Options` header.
	// Optional. Default value "nosniff".
	ContentTypeOptions string
	// Sets X-Frame-Options header that can be used to indicate whether or not a browser should be allowed to render a
	// page in a <frame>, <iframe>, <embed> or <object>. Valid values are deny, sameorigin and allow-from uri. Default
	// is deny
	FrameOptions string
	// Sets the `Strict-Transport-Security` header to indicate how
	// long (in seconds) browsers should remember that this site is only to
	// be accessed using HTTPS. Default is 0
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Strict-Transport-Security
	// HSTSMaxAge int
	// When enabled this rule applies to all of the site's subdomains as well to Strict-Transport-Security
	// HSTSIncludeSubdomains bool

	// ReferrerPolicy sets the `Referrer-Policy` header providing security against
	// leaking potentially sensitive request paths to third parties.
	// Optional. Default value "".  https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Referrer-Policy
	ReferrerPolicy string
}

// SetDefaultHeaders will set certain default headers specified by the user
func (defaultHeaders *DefaultHeaders) SetDefaultHeaders(h router.Handle) router.Handle {
	return func(w http.ResponseWriter, request *http.Request, params map[string]string) {
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
		h(w, request, params)

	}
}

// SetSecureHeaders sets some default security headers
func (secureHeaders SecureHeaders) SetSecureHeaders(h router.Handle) router.Handle {
	return func(w http.ResponseWriter, request *http.Request, params map[string]string) {
		if secureHeaders.XSSProtection != "" {
			w.Header().Set("X-XSS-Protection", secureHeaders.XSSProtection)
		} else {

			w.Header().Set("X-XSS-Protection", "1; mode=block")
		}

		if secureHeaders.ContentTypeOptions != "" {
			w.Header().Set("X-Content-Type-Options", secureHeaders.ContentTypeOptions)
		} else {
			w.Header().Set("X-Content-Type-Options", "nosniff")
		}

		if secureHeaders.FrameOptions != "" {
			w.Header().Set("X-Frame-Options", secureHeaders.FrameOptions)
		} else {
			w.Header().Set("X-Frame-Options", "deny")
		}

		w.Header().Set("Referrer-Policy", secureHeaders.ReferrerPolicy)
		h(w, request, params)
	}
}
