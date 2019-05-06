package response

import (
	"encoding/json"
	"net/http"
)

func SuccessResponse(status int, responseWriter http.ResponseWriter, body interface{}) {
	responseWriter.WriteHeader(status)

	bytesRep, _ := json.Marshal(body)
	_, _ = responseWriter.Write(bytesRep)

}
