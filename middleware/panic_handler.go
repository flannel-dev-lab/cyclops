package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/flannel-dev-lab/cyclops/v2/logger"
)

// PanicHandler takes care of recovering from panic if any unforeseen error occurs in the execution logic and makes sure
// that the server does not stop
func PanicHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		ctx := request.Context()

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

				ctx = logger.AddKey(ctx, "stack-trace", string(debug.Stack()))
				ctx = logger.AddKey(ctx, "err", err.Error())

				errData, _ := json.Marshal(map[string]string{"error": err.Error()})

				http.Error(responseWriter, string(errData), http.StatusInternalServerError)
				return
			}
		}()
		h.ServeHTTP(responseWriter, request)
	}
}
