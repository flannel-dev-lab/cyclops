package middleware

import (
	"errors"
	"net/http"
)

func PanicHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
		var err error
		defer func() {
			r := recover()

			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					err = errors.New("unknown error")
				}
				http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
				return
				// TODO respond error
				// middlewareService.respondFailure(http.StatusUnprocessableEntity, err, string(debug.Stack()), responseWriter)
			}
		}()
		h.ServeHTTP(responseWriter, request)
	})
}
