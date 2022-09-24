// Package input deals with getting the values from form and query parameters
package input

import (
	"mime/multipart"
	"net/http"
)

// Query retrieves the url parameters from the request
func Query(key string, request *http.Request) string {
	return request.URL.Query().Get(key)
}

// Form retrieves the form parameters from the request
func Form(key string, request *http.Request) string {
	return request.FormValue(key)
}

// FileContent retrieves the file contents from the requests, else returns error
func FileContent(r *http.Request, key string) (multipart.File, *multipart.FileHeader, error) {
	return r.FormFile(key)
}
