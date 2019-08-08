package response

import "net/http"

func Redirect(w http.ResponseWriter, r *http.Request, url string, statusCode int) {
	http.Redirect(w, r, url, statusCode)
}
