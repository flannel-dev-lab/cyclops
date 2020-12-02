package response

import (
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

	request, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
	if err != nil {
		t.Error(err)
	}

	request.Header.Set("Content-Type", "application/json")

	request.URL.Query().Set("test", "test")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Error("response codes do not match for error response")
	}
}
