// input deals with getting the values from form and query parameters
package input

import "net/http"

// Get retrieves the url parameters from the request
func Query(key string, request *http.Request) string {
	return request.URL.Query().Get(key)
}

// Get retrieves the form parameters from the request
func Form(key string, request *http.Request) string {
	return request.FormValue(key)
}
