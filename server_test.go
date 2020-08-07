package microframework

import (
	"net/http"
	"testing"
)

var expected string = "AB"

func TestNewServer(t *testing.T) {
	c := &Config{Address: ":9001"}
	s := NewServer(c)

	if len(s.GlobalMiddleware) != 0 {
		t.Errorf(
			"Server contains global middleware when there shouldn't be any: %v",
			s.GlobalMiddleware)
	}
}

func TestGlobalMiddleware(t *testing.T) {
	c := &Config{Address: ":9001"}
	s := NewServer(c, middlewareA)

	if l := len(s.GlobalMiddleware); l != 1 {
		t.Errorf("s.GlobalMiddleware: expected: 1, actual: %d", l)
	}

	s.AddGlobalMiddleware(middlewareA)

	if l := len(s.GlobalMiddleware); l != 2 {
		t.Errorf("s.GlobalMiddleware: expected: 2, actual: %d", l)
	}
}

func middlewareA(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		return next(w, r)
	}
}
