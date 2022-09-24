// Package cookie deals with set, get and deleting of HTTP Cookies
package cookie

import (
	"errors"
	"net/http"
	"time"
)

const (
	// DefaultExpiry is the default expiry for a cookie if a user does not set it
	DefaultExpiry = 3600
)

// CyclopsCookie is a struct to hold cookie info
type CyclopsCookie struct {
	// Name of cookie
	Name string
	// Value of cookie
	Value string
	// The Domain and Path attributes define the scope of the cookie. They essentially tell the browser what website
	// the cookie belongs to
	Path   string
	Domain string
	// Secure attribute is meant to keep cookie communication limited to encrypted transmission, directing browsers
	// to use cookies only via secure/encrypted connections
	Secure bool
	// HttpOnly attribute directs browsers not to expose cookies through channels other than HTTP (and HTTPS) requests.
	// This means that the cookie cannot be accessed via client-side scripting languages (notably JavaScript),
	// and therefore cannot be stolen easily via cross-site scripting
	HttpOnly bool
	// SameSite when enabled, a cookie can only be sent in requests originating from the same origin as the target
	// domain, it helps to prevent XSRF attacks
	SameSite http.SameSite
	// Expires specifies time in seconds for a cookie to expire
	Expires time.Duration
	MaxAge  int
}

// SetCookie  is used to set a cookie to the responseWriter
func (cyclopsCookie CyclopsCookie) SetCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{}
	cookie.Name = cyclopsCookie.Name
	cookie.Value = cyclopsCookie.Value
	if cyclopsCookie.Path == "" {
		cookie.Path = "/"
	} else {
		cookie.Path = cyclopsCookie.Path
	}
	cookie.Domain = cyclopsCookie.Domain
	cookie.Secure = cyclopsCookie.Secure
	cookie.HttpOnly = cyclopsCookie.HttpOnly

	if cyclopsCookie.Expires == 0 {
		cookie.Expires = time.Now().Add(DefaultExpiry * time.Second)
	} else {
		cookie.Expires = time.Now().Add(cyclopsCookie.Expires * time.Second)
	}

	if cyclopsCookie.SameSite == 0 {
		cookie.SameSite = http.SameSiteNoneMode
	} else {
		cookie.SameSite = cyclopsCookie.SameSite
	}

	http.SetCookie(w, cookie)
}

// GetCookie retrieves a cookie based on the cookie name from request, returns error if cookie does not exist
func (cyclopsCookie CyclopsCookie) GetCookie(r *http.Request, name string) (*http.Cookie, error) {
	if r.Method == http.MethodTrace {
		return nil, errors.New("method not allowed")
	}
	return r.Cookie(name)
}

// GetAll returns array of HTTP cookies
func (cyclopsCookie CyclopsCookie) GetAll(r *http.Request) []*http.Cookie {
	return r.Cookies()
}

// Delete is used to delete a cookie from the browser by setting the expires attribute to 0
func (cyclopsCookie CyclopsCookie) Delete(w http.ResponseWriter, cookie *http.Cookie) {
	cookie.Expires = time.Now().Add(0 * time.Second)

	http.SetCookie(w, cookie)
}
