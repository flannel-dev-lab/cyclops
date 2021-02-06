package response

import (
	"encoding/json"
	"net/http"
)

// Error holds the error in string
type Error struct {
	Error string `json:"error"`
}

// ErrorResponse handles the logging and structuring of sending error to the user
func ErrorResponse(status int, message string, responseWriter http.ResponseWriter) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(status)

	response := &Error{Error: message}

	bytesRep, _ := json.Marshal(response)
	_, _ = responseWriter.Write(bytesRep)
}
