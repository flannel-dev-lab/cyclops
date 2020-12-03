package input

import (
	"bytes"
	"fmt"
	"github.com/flannel-dev-lab/cyclops/v2/router"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestQuery(t *testing.T) {
	t.Parallel()

	request, err := http.NewRequest(http.MethodPost, "http://localhost", nil)
	if err != nil {
		t.Error(err)
	}

	q := request.URL.Query()
	q.Add("test", "value")

	request.URL.RawQuery = q.Encode()

	if Query("test", request) != "value" {
		t.Error("query values do not match")
	}

	if Query("empty", request) != "" {
		t.Error("non defined key has a value set")
	}
}

func Test_Form(t *testing.T) {
	t.Parallel()

	r := router.New(true, nil, nil)
	r.Post("/use", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.ParseForm())
		fmt.Println(r.Form.Get("test"))
		fmt.Println(Form("test", r))
	})

	data := url.Values{}
	data.Set("client_id", "test")

	req, _ := http.NewRequest("POST", "/use", bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	fmt.Println("OUT", Form("test", req))
}

// Benchmarks
func BenchmarkQuery(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()

	request, err := http.NewRequest(http.MethodPost, "http://localhost", nil)
	if err != nil {
		b.Error(err)
	}

	q := request.URL.Query()
	q.Add("test", "value")

	request.URL.RawQuery = q.Encode()

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		Query("test", request)
	}
}
