package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// LogObject contains keys for the logs
type LogObject struct {
	// Current timestamp of request
	Timestamp     string `json:"timestamp"`
	// RemoteAddress contains the IP of the server/ the IP address of the proxy
	RemoteAddress string `json:"remote_address"`
	// TrueIP contains the IP of the original requester
	TrueIP string `json:"true_ip"`
	// Method contains the http method requested
	Method        string `json:"method"`
	// Path contains the http path requested
	Path          string `json:"path"`
	// Host contains the IP of host
	Host          string `json:"host"`
	// Protocol contains http version
	Protocol      string `json:"protocol"`
}

// logWriter struct that implements Write for logger
type logWriter struct {
}

// Write prints logs to stdout in JSON
func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Printf("%s\n", string(bytes))
}

// RequestLogger intercepts logs from the requests and prints them to stdout
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
