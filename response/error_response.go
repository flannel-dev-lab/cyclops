package response

import (
	"encoding/json"
	"github.com/flannel-dev-lab/cyclops/alerts"
	"log"
	"net/http"
)

type Error struct {
	Error string `json:"error"`
}

func ErrorResponse(status int, err error, message string, responseWriter http.ResponseWriter, sendAlert bool, alert alerts.Alert) {
	responseWriter.WriteHeader(status)

	if sendAlert {
		alert.CaptureError(
			err,
			err.Error())
	}

	log.Printf("%v", err)

	response := &Error{Error: message}

	bytesRep, _ := json.Marshal(response)
	_, err = responseWriter.Write(bytesRep)

}
