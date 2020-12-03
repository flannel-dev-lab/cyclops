package middleware

import (
	"github.com/flannel-dev-lab/cyclops/v2/router"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS_CORSHandler(t *testing.T) {
	cors := []CORS{
		{"*", true, []string{"Content-Type"}, []string{"HEAD"}, []string{"Content-Type"}, 300},
	}

	for _, testCase := range cors {
		r := router.New(true, nil, nil)
		r.Post("/use", NewChain(testCase.CORSHandler).Then(func(w http.ResponseWriter, r *http.Request) {}))
		r.Options("/use", NewChain(testCase.CORSHandler).Then(func(w http.ResponseWriter, r *http.Request) {}))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/use", nil)

		r.ServeHTTP(w, req)

		if w.Header().Get("Access-Control-Allow-Origin") != testCase.AllowedOrigin {
			t.Errorf("expected %s got %s", testCase.AllowedOrigin, w.Header().Get("Access-Control-Allow-Origin"))
		}

		w = httptest.NewRecorder()
		req, _ = http.NewRequest("OPTIONS", "/use", nil)

		r.ServeHTTP(w, req)

		if w.Header().Get("Access-Control-Allow-Origin") != testCase.AllowedOrigin {
			t.Errorf("expected %s got %s", testCase.AllowedOrigin, w.Header().Get("Access-Control-Allow-Origin"))
		}
	}
}
