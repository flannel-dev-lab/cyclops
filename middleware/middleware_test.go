package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestContentTypeFilter_JSON(t *testing.T) {
	testServer := httptest.NewServer(RequestContentTypeFilter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Root!\n")
	}), "application/json"))

	request, err := http.NewRequest("POST", testServer.URL, nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if response.StatusCode != 200 {
		t.Fatal("json response gave non 200")
	}
}

func TestRequestContentTypeFilter_BadContentType(t *testing.T) {
	testServer := httptest.NewServer(RequestContentTypeFilter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Root!\n")
	}), "application/json"))

	request, err := http.NewRequest("POST", testServer.URL, nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	request.Header.Set("Content-Type", "application/data")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if response.StatusCode == 200 {
		t.Fatal("bad content gave 200")
	}
}

func TestSetHeaders(t *testing.T) {
	testServer := httptest.NewServer(SetSecureHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Root!\n")
	})))

	request, err := http.NewRequest("POST", testServer.URL, nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if response.Header.Get("X-Frame-Options") != "deny" {
		t.Fatal("unable to set headers")
	}
}

// TODO Host Header validator
// TODO JWT
// TODO CORS
