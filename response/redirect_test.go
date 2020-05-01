package response

import (
	"fmt"
	"github.com/flannel-dev-lab/cyclops/requester"
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

	response, err := requester.Get(testServer.URL, map[string]string{"Content-Type": "application/json"}, map[string]string{"test": "test"})
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