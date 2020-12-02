package response

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRedirect(t *testing.T) {
	redirectServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Root!")
	}))

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Redirect(w, r, redirectServer.URL, http.StatusFound)
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

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
	}

	bodyString := string(bodyBytes)

	if bodyString != "Root!" {
		t.Error("response body do not match")
	}
}
