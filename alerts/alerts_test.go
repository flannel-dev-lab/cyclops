package alerts

import "testing"

func TestInitiateAlerting_SupportedSystem(t *testing.T) {
	_, err := InitiateAlerting("sentry", "http://sentry.io", "dev", true)

	if err != nil {
		t.Error("supporting system gave an error")
	}
}

func TestInitiateAlerting_UnsupportedSystem(t *testing.T) {
	_, err := InitiateAlerting("borg", "http://borg.io", "dev", true)

	if err == nil {
		t.Error("unsupporting system gave no error")
	}
}
