package router

import (
	"net/http"
	"testing"
)

func TestValidMethod(t *testing.T) {
	cases := []struct{
		Method string
		Expected bool
	}{
		{http.MethodGet, true},
		{"FUNNY", false},
	}

	for _, test := range cases {
		value := validMethod(test.Method)
		if value != test.Expected {
			t.Errorf("Expected %t Got %t", test.Expected, value)
		}
	}
}