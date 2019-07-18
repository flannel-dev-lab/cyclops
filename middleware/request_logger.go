package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type LogObject struct {
	Timestamp     string `json:"timestamp"`
	RemoteAddress string `json:"remote_address"`
	Method        string `json:"method"`
	Path          string `json:"path"`
	Host          string `json:"host"`
	Protocol      string `json:"protocol"`
}

func (logObject LogObject) Write(bytes []byte) (int, error) {
	logData, _ := json.Marshal(logObject)
	return fmt.Print(string(logData))
}

// Middleware to log access logs
func RequestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {

		logObject := LogObject{
			Timestamp:     time.Now().UTC().Format("2006-01-02T15:04:05.999Z"),
			RemoteAddress: request.RemoteAddr,
			Method:        request.Method,
			Path:          request.URL.Path,
			Host:          request.Host,
			Protocol:      request.Proto,
		}

		log.SetOutput(logObject)
		defer log.Println()
		h.ServeHTTP(w, request)
	})
}
