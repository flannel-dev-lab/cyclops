// cookie deals with setting and getting cookies
package cookie

import (
	"errors"
	"net/http"
	"time"
)

const DefaultExpiry = 3600

// CyclopsCookie is an object to hold cookie info
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
	// StrictSameSite helps to prevent XSRF attacks
	StrictSameSite bool
	// Expires specifies time in seconds for a cookie to expire
	Expires time.Duration
}

// SetCookie Sets a cookie to the responseWriter
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

	if cyclopsCookie.StrictSameSite {
		cookie.SameSite = http.SameSiteStrictMode
	}

	http.SetCookie(w, cookie)
}

// GetCookie retrieves a cookie based on the key from request
func (cyclopsCookie CyclopsCookie) GetCookie(r *http.Request, key string) (*http.Cookie, error) {
	if r.Method == http.MethodTrace {
		return nil, errors.New("method not allowed")
	}
	return r.Cookie(key)
}

// GetAll returns array of cookies
func (cyclopsCookie CyclopsCookie) GetAll(r *http.Request) []*http.Cookie {
	return r.Cookies()
}
