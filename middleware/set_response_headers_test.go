package middleware

import (
	"fmt"
	"github.com/flannel-dev-lab/cyclops/v2/router"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSecureHeaders_SetSecureHeaders(t *testing.T) {
	cases := []SecureHeaders{
		{"1; mode=block", "nosniff", "sameorigin", ""},
		{"", "nosniff", "sameorigin", ""},
		{"1; mode=block", "", "sameorigin", ""},
		{"1; mode=block", "nosniff", "", ""},
	}

	for _, testCase := range cases {

		r := router.New(true, nil, nil)
		r.Post("/use", NewChain(testCase.SetSecureHeaders).Then(func(w http.ResponseWriter, r *http.Request) {}))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/use", nil)

		r.ServeHTTP(w, req)
		fmt.Println(w.Header())

		if w.Header().Get("X-Content-Type-Options") != "nosniff" {
			t.Errorf("%s: headers do not match expected '%s' got '%s'", t.Name(), testCase.ContentTypeOptions, w.Header().Get("X-Content-Type-Options"))
		}

		if w.Code != http.StatusOK {
			t.Errorf("%s: expected 200 got %d", t.Name(), w.Code)
		}

	}
}

func TestDefaultHeaders_SetDefaultHeaders(t *testing.T) {
	cases := []DefaultHeaders{
		{"test", []string{"en"}, "application/json", "Wed, 21 Oct 2015 07:28:00 GMT", "Mon"},
		{"", []string{"en"}, "application/json", "Wed, 21 Oct 2015 07:28:00 GMT", "Mon"},
		{"test", []string{"en"}, "", "Wed, 21 Oct 2015 07:28:00 GMT", "Mon"},
	}

	for _, testCase := range cases {

		r := router.New(true, nil, nil)
		r.Post("/use", NewChain(testCase.SetDefaultHeaders).Then(func(w http.ResponseWriter, r *http.Request) {}))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/use", nil)

		r.ServeHTTP(w, req)

		if w.Header().Get("Content-Type") != testCase.ContentType {
			t.Errorf("%s: headers do not match expected '%s' got '%s'", t.Name(), testCase.ContentType, w.Header().Get("Content-Type"))
		}

		if w.Code != http.StatusOK {
			t.Errorf("%s: expected 200 got %d", t.Name(), w.Code)
		}

	}
}
