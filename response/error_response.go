package response

import (
	"encoding/json"
	"fmt"
	"github.com/flannel-dev-lab/cyclops/alerts"
	"log"
	"net/http"
	"time"
)

type Error struct {
	Error string `json:"error"`
}

type ErrorLogger struct {
	Timestamp   string `json:"timestamp"`
	Error       string `json:"error"`
	StatusCode  int    `json:"status_code"`
	AlertStatus bool   `json:"alert_status"`
}

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
