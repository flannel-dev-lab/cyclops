// input deals with getting the values from form and query parameters
package input

import (
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

// Query retrieves the url parameters from the request
func Query(key string, request *http.Request) string {
	return request.URL.Query().Get(key)
}

// TrimmedParamNames - Gets aaa url parameter names list without the leading :
func TrimmedParamNames(r *http.Request) []string {
	var names []string
	for k := range r.URL.Query() {
		if strings.HasPrefix(k, ":") {
			names = append(names, strings.TrimPrefix(k, ":"))
		}
	}
	return names
}

// AddParam is useful for middlewares
// Appends :name=value onto a blank request query string or appends &:name=value
// onto a non-blank request query string
func AddParam(r *http.Request, name, value string) {
	q := url.QueryEscape(":"+name) + "=" + url.QueryEscape(value)
	if r.URL.RawQuery != "" {
		r.URL.RawQuery += "&" + q
	} else {
		r.URL.RawQuery += q
	}
}

// Form retrieves the form parameters from the request
func Form(key string, request *http.Request) string {
	return request.FormValue(key)
}

func FileContent(r *http.Request, key string) (multipart.File, *multipart.FileHeader, error) {
	return r.FormFile(key)
}
