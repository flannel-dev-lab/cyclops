package input

import (
	"net/http"
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

func TestTrimmedParamNames(t *testing.T) {
	t.Parallel()

	request, err := http.NewRequest(http.MethodPost, "http://localhost", nil)
	if err != nil {
		t.Error(err)
	}

	q := request.URL.Query()
	q.Add(":test", "value")
	q.Add("test2", "value2")

	request.URL.RawQuery = q.Encode()

	if len(TrimmedParamNames(request)) != 1 {
		t.Error("existing trimmed param not present in result")
	}
}

func TestAddParam(t *testing.T) {
	t.Parallel()

	request, err := http.NewRequest(http.MethodPost, "http://localhost", nil)
	if err != nil {
		t.Error(err)
	}

	AddParam(request, "test", "value")
	AddParam(request, ":test", "value2")

	if len(TrimmedParamNames(request)) != 2 {
		t.Error("no parameters were added")
	}
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

func BenchmarkTrimmedParamNames(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()

	request, err := http.NewRequest(http.MethodPost, "http://localhost", nil)
	if err != nil {
		b.Error(err)
	}

	q := request.URL.Query()
	q.Add(":test", "value")
	q.Add("test2", "value2")

	request.URL.RawQuery = q.Encode()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		TrimmedParamNames(request)
	}
}

func BenchmarkAddParam(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()

	request, err := http.NewRequest(http.MethodPost, "http://localhost", nil)
	if err != nil {
		b.Error(err)
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		AddParam(request, "test", "value")
	}
}
