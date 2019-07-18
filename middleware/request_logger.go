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

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Printf("%s\n", string(bytes))
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
		log.SetFlags(0)
		log.SetOutput(new(logWriter))

		logData, _ := json.Marshal(logObject)

		defer log.Println(fmt.Sprintf("%s", logData))
		h.ServeHTTP(w, request)
	})
}
