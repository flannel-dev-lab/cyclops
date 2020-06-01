package cookie

import (
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
	cookies := []CyclopsCookie{
		{
			Name:     "name",
			Value:    "test",
			Domain:   "google.com",
			Secure:   false,
			HttpOnly: false,
			SameSite: http.SameSiteNoneMode,
			Expires:  3600,
			MaxAge:   3600,
		},
		{
			Name:     "name",
			Value:    "test",
			Path:     "/",
			Domain:   "google.com",
			Secure:   false,
			HttpOnly: false,
			SameSite: http.SameSiteNoneMode,
			Expires:  0,
			MaxAge:   0,
		},
		{
			Name:     "name",
			Value:    "test",
			Path:     "/",
			Domain:   "google.com",
			Secure:   false,
			HttpOnly: false,
			SameSite: 0,
			Expires:  3600,
			MaxAge:   3600,
		},
		{
			Name:     "name",
			Value:    "test",
			Path:     "/",
			Domain:   "",
			Secure:   false,
			HttpOnly: false,
			SameSite: http.SameSiteNoneMode,
			Expires:  3600,
			MaxAge:   3600,
		},
	}

	for _, cookie := range cookies {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie.SetCookie(w)
			w.WriteHeader(200)
			return
		}))

		request, _ := http.NewRequest("GET", srv.URL, nil)
		client := &http.Client{}

		response, err := client.Do(request)
		if err != nil {
			t.Fatalf("%v", err)
		}

		if response.Cookies()[0].Value != "test" {
			t.Fatalf("%s", "cookie values do not match")
		}
	}
}
