package requester

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGet(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Root!\n")
	}))

	response, err := Get(testServer.URL, map[string]string{"Content-Type": "application/json"}, map[string]string{"test": "test"})
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Errorf("unexpected status code %d", response.StatusCode)
	}
}


func TestPost(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Root!\n")
	}))

	response, err := Post(testServer.URL, map[string]string{"Content-Type": "application/json"}, map[string]string{"test": "test"}, nil)
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Errorf("unexpected status code %d", response.StatusCode)
	}
}

func TestDelete(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Root!\n")
	}))

	response, err := Delete(testServer.URL, map[string]string{"Content-Type": "application/json"}, map[string]string{"test": "test"}, nil)
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Errorf("unexpected status code %d", response.StatusCode)
	}
}

func TestPut(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Root!\n")
	}))

	response, err := Put(testServer.URL, map[string]string{"Content-Type": "application/json"}, map[string]string{"test": "test"}, nil)
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Errorf("unexpected status code %d", response.StatusCode)
	}
}

func TestPatch(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Root!\n")
	}))

	response, err := Patch(testServer.URL, map[string]string{"Content-Type": "application/json"}, map[string]string{"test": "test"}, nil)
	if err != nil {
		t.Error(err)
	}

	if response.StatusCode != 200 {
		t.Errorf("unexpected status code %d", response.StatusCode)
	}
}