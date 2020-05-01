package sentry

import (
	"errors"
	"os"
	"testing"
)

var (
	sentry Sentry
)

func TestMain(m *testing.M) {
	sentry.Enabled = true
	sentry.Environment = "test"
	sentry.DSN = "https://sentry.com"

	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestSentry_Bootstrap(t *testing.T) {
	if err := sentry.Bootstrap(); err == nil {
		t.Error("Bad sentry err gave no err")
	}
}

func TestSentry_CaptureError(t *testing.T) {
	sentry.CaptureError(errors.New("test-error"), "test-error")
}

func BenchmarkSentry_Bootstrap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sentry.Bootstrap()
	}
}

func BenchmarkSentry_CaptureError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sentry.CaptureError(errors.New("test-error"), "test-error")
	}
}
