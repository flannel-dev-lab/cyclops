package cookie

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestCyclopsCookie_SetCookie(t *testing.T) {
	var cyclopsCookie CyclopsCookie
	cyclopsCookie.Name = "Hello"
	cyclopsCookie.Value = "Cyclops"
	cyclopsCookie.Path = "/"

	_, _ = http.NewRequest("GET", "https://www.google.com", nil)
	response := httptest.NewRecorder()
	cyclopsCookie.SetCookie(response)
}

func TestCyclopsCookie_GetAll(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var cyclopsCookie CyclopsCookie
		cyclopsCookie.Name = "Hello"
		cyclopsCookie.Value = "Cyclops"
		cyclopsCookie.SetCookie(w)
		w.WriteHeader(200)
		return
	}))

	request, _ := http.NewRequest("GET", srv.URL, nil)
	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		t.Errorf("%v", err)
	}

	cookies := CyclopsCookie{}.GetAll(response.Request)
	fmt.Println(cookies, err)

	if response.Cookies()[0].Value != "Cyclops" {
		t.Error("cookie values does not match")
	}
}
