package middleware

import (
	"errors"
	"github.com/flannel-dev-lab/cyclops/router"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPanicHandler(t *testing.T) {
	cases := []interface{}{
		errors.New("error panic"),
		1,
		"string panic",
	}

	for _, testCase := range cases {
		r := router.New(true, nil, nil)
		r.Post("/use", NewChain(PanicHandler).Then(func(w http.ResponseWriter, r *http.Request) {
			panic(testCase)
		}))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/use", nil)

		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("%s: expected 500 got %d", t.Name(), w.Code)
		}
	}

}
