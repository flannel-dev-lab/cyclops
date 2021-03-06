package response

import (
	"encoding/json"
	"net/http"
)

// SuccessResponse sends a success response to the user
func SuccessResponse(status int, responseWriter http.ResponseWriter, body interface{}) {
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(status)

	bytesRep, _ := json.Marshal(body)
	_, _ = responseWriter.Write(bytesRep)

}
