package microframework

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestQueue(t *testing.T) {
	middleware := []Middleware{decodeChar}
	q := newQueue(controller, middleware, func(e error) {
		t.Errorf("%v", e)
	}, nil)

	// create a test server with the queue
	ts := httptest.NewServer(q)
	defer ts.Close()

	reader := bytes.NewReader([]byte(`{"name": "Shao Khan", "wins": 9}`))

	// use the test server's URL as the address
	req, e := http.NewRequest(http.MethodPost, ts.URL, reader)

	if e != nil {
		t.Fatal(e)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	// send the request to the test server
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
	errorLogger := func(e error) {
		// ensure the error logger is called and receives an HttpError
		_, ok := e.(HttpError)

		if !ok {
			t.Errorf("Expected an HttpError, received: %v", e)
		}
	}

	q := queue{el: errorLogger}

	q.handleError(recorder, e)

	res := recorder.Result()

	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("%v returned as error", res)
	}
}

type Character struct {
	name string
	wins float64
}

func (c *Character) toContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "char", c)
}

func charFromContext(ctx context.Context) (*Character, error) {
	v := ctx.Value("char")
	c, ok := v.(*Character)

	if !ok {
		return nil, fmt.Errorf("Could not assert %v as *Character", v)
	}

	return c, nil
}

func charFromJSON(bytes []byte) (*Character, error) {
	c := &Character{}
	e := json.Unmarshal(bytes, c)

	return c, e
}

func decodeChar(r *http.Request) error {
	defer r.Body.Close()
	bytes, e := ioutil.ReadAll(r.Body)

	if e != nil {
		return e
	}

	c, e := charFromJSON(bytes)

	if e != nil {
		return e
	}

	*r = *r.WithContext(c.toContext(r.Context()))

	return nil
}

func controller(w http.ResponseWriter, r *http.Request) error {
	_, e := charFromContext(r.Context())

	if e != nil {
		return e
	}

	w.WriteHeader(http.StatusOK)

	return nil
}
