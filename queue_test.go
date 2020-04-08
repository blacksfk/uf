package microframework

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQueue(t *testing.T) {
	middleware := []Middleware{decodeChar}
	q := newQueue(controller, middleware)
	ts := httptest.NewServer(q)
	defer ts.Close()

	reader := bytes.NewReader([]byte(`{"name": "Shao Khan", "wins": 9}`))
	req, e := http.NewRequest(http.MethodPost, ts.URL, reader)

	if e != nil {
		t.Fatal(e)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, e := client.Do(req)

	if e != nil {
		t.Fatal(e)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("%v", res)
	}
}

func TestHandleError(t *testing.T) {
	recorder := httptest.NewRecorder()
	e := errors.New("You broke the gearbox")

	handleError(recorder, e)

	res := recorder.Result()

	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("%v return as error", res)
	}
}

type Character struct {
	name string
	wins float64
}

func (c *Character) toContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "char", c)
}

func charFromContext(ctx context.Context) *Character {
	return ctx.Value("char").(*Character)
}

func charFromJSON(bytes []byte) (*Character, error) {
	c := &Character{}
	e := json.Unmarshal(bytes, c)

	return c, e
}

func decodeChar(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		defer r.Body.Close()
		bytes, e := ioutil.ReadAll(r.Body)

		if e != nil {
			return e
		}

		c, e := charFromJSON(bytes)

		if e != nil {
			return e
		}

		r = r.WithContext(c.toContext(r.Context()))

		return next(w, r)
	}
}

func controller(w http.ResponseWriter, r *http.Request) error {
	c := charFromContext(r.Context())

	if c == nil {
		return InternalServerError("Character not in context")
	}

	w.WriteHeader(http.StatusOK)

	return nil
}
