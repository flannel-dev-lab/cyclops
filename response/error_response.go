package response

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Error holds the error in string
type Error struct {
	Error string `json:"error"`
}

// ErrorLogger contains fields to print error log to stdout
type ErrorLogger struct {
	// Current timestamp of request
	Timestamp string `json:"timestamp"`
	// Error in string format
	Error string `json:"error"`
	// Response status code
	StatusCode int `json:"status_code"`
}

// ErrorResponse handles the logging and structuring of sending error to the user
func ErrorResponse(status int, err error, message string, responseWriter http.ResponseWriter) {
	responseWriter.WriteHeader(status)

	errorLog := ErrorLogger{
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Error:       err.Error(),
		StatusCode:  status}
	logData, _ := json.Marshal(errorLog)

	log.Printf(fmt.Sprintf("%s", logData))

	response := &Error{Error: message}

	bytesRep, _ := json.Marshal(response)
	_, err = responseWriter.Write(bytesRep)
}
