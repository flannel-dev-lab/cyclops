package response

import (
	"github.com/flannel-dev-lab/cyclops/requester"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSuccessResponse(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SuccessResponse(
			http.StatusOK,
			w,
			nil)
		return
	}))

	response, err := requester.Get(testServer.URL, map[string]string{"Content-Type": "application/json"}, map[string]string{"test": "test"})
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Error("response codes do not match for error response")
	}
}
