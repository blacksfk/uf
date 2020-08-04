package microframework

import (
	"net/http"
	"testing"
)

func TestNewGroup(t *testing.T) {
	server := NewServer(&Config{Address: ":9001"})
	group := server.NewGroup("/nothing", doNothing)

	if len(group.middleware) != 1 {
		t.Fatalf("middleware not appended: %+v", group)
	}
}

func TestGroupMiddleware(t *testing.T) {
	server := NewServer(&Config{Address: ":9001"})
	group := server.NewGroup("/nothing", doNothing).Middleware(doNothing)

	if len(group.middleware) != 2 {
		t.Fatalf("middleware not appended: %+v", group)
	}
}

func TestGroupMethods(t *testing.T) {
	server := NewServer(&Config{Address: ":9001"})
	group := server.NewGroup("/nothing", doNothing)

	group.Get(handleNothing, doNothing).Post(handleNothing).Middleware(doNothing, doNothing).Put(handleNothing).Patch(handleNothing).Delete(handleNothing)

	if len(group.middleware) != 3 {
		t.Fatalf("middleware not appended: %+v", group)
	}
}

func doNothing(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		return next(w, r)
	}
}

func handleNothing(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(200)

	return nil
}
