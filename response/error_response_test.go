package response

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestErrorResponse(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ErrorResponse(
			http.StatusBadRequest,
			"test error",
			w)
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

	if response.StatusCode != http.StatusBadRequest {
		t.Error("response codes do not match for error response")
	}
}

type MockAlert struct {
}

func (m MockAlert) Bootstrap() error {
	return nil
}

func (m MockAlert) CaptureError(err error, message string) {}

func TestErrorResponse_SendAlert(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ErrorResponse(
			http.StatusBadRequest,
			"test error",
			w)
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

	if response.StatusCode != http.StatusBadRequest {
		t.Error("response codes do not match for error response")
	}
}
