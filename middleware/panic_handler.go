package middleware

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"runtime/debug"
)

// PanicHandler takes care of recovering from panic if any unforseen error occurs in the execution logic and makes sure
// that the server does not stop
func PanicHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
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
				log.Print(string(debug.Stack()))
				log.Println(err.Error())

				errData, _ := json.Marshal(map[string]string{"error": err.Error()})

				http.Error(responseWriter, string(errData), http.StatusInternalServerError)
				return
			}
		}()
		h.ServeHTTP(responseWriter, request)
	}
}
