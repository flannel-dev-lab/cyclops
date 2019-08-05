package response

import (
	"encoding/json"
	"fmt"
	"github.com/flannel-dev-lab/cyclops/alerts"
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
	// Tells if an alert is sent or not
	AlertStatus bool `json:"alert_status"`
}

// ErrorResponse handles the logging and structuring of sending error to the user
func ErrorResponse(status int, err error, message string, responseWriter http.ResponseWriter, sendAlert bool, alert alerts.Alert) {
	responseWriter.WriteHeader(status)

	if sendAlert {
		alert.CaptureError(
			err,
			err.Error())
	}

	errorLog := ErrorLogger{
		Timestamp:   time.Now().UTC().Format("2006-01-02T15:04:05.999Z"),
		Error:       err.Error(),
		StatusCode:  status,
		AlertStatus: sendAlert}
	logData, _ := json.Marshal(errorLog)

	log.Printf(fmt.Sprintf("%s", logData))

	response := &Error{Error: message}

	bytesRep, _ := json.Marshal(response)
	_, err = responseWriter.Write(bytesRep)
}
