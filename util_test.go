package microframework

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type GT1 struct {
	Manufacturer, Model string
	Debut               int
}

func TestSendJSON(t *testing.T) {
	// create a new test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// try to send back JSON
		nissan := GT1{"Nissan", "R390 LM", 1997}
		e := SendJSON(w, nissan)

		if e != nil {
			// something went wrong with the encoding?
			t.Fatal(e)
		}
	}))

	defer ts.Close()
	// send a request
	res, e := http.Get(ts.URL)

	if e != nil {
		t.Fatal(e)
	}

	defer res.Body.Close()
	// check the correct content type was set and decode the body
	checkContentType(t, res)
	b, e := io.ReadAll(res.Body)

	if e != nil {
		t.Fatal(e)
	}

	// ensure the data properly unmarshals
	nissan := GT1{}
	e = json.Unmarshal(b, &nissan)

	if e != nil {
		t.Error(e)
	}
}

func TestSendErrorJSON(t *testing.T) {
	// create a test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// send back an HttpError as JSON
		he := BadRequest("You sent garbage")
		e := SendErrorJSON(w, he)

		if e != nil {
			t.Fatal(e)
		}
	}))

	defer ts.Close()
	// send the request
	res, e := http.Get(ts.URL)

	if e != nil {
		t.Fatal(e)
	}

	defer res.Body.Close()
	// ensure the correct content type and decode the body
	checkContentType(t, res)
	b, e := io.ReadAll(res.Body)

	if e != nil {
		t.Fatal(e)
	}

	var m map[string]interface{}
	e = json.Unmarshal(b, &m)

	if e != nil {
		t.Fatal(e)
	}

	_, ok := m["code"].(float64)

	if !ok {
		t.Fatalf("Could not assert %v as float64", m["code"])
	}
}

func TestReadBody(t *testing.T) {
	b, e := json.Marshal(GT1{"Toyota", "TS020", 1998})

	if e != nil {
		t.Fatal(e)
	}

	reader := bytes.NewReader(b)
	r, e := http.NewRequest(http.MethodPost, "http://example.com", reader)

	if e != nil {
		t.Fatal(e)
	}

	r.Header.Set("Content-Type", "application/json")

	// accept any content-type header
	_, e = ReadBody(r)

	if e != nil {
		t.Error(e)
	}

	r, e = http.NewRequest(http.MethodPost, "http://example.com", reader)

	if e != nil {
		t.Fatal(e)
	}

	r.Header.Set("Content-Type", "text/plain")

	// only allow json
	_, e = ReadBody(r, "application/json")

	if e == nil {
		t.Error("Allowed text/plain when expecting application/json")
	}

	r, e = http.NewRequest(http.MethodPost, "http://example.com", reader)

	r.Header.Set("Content-Type", "application/json")

	if e != nil {
		t.Fatal(e)
	}

	// allow both json and x-www-form-urlencoded
	_, e = ReadBody(r, "application/json", "application/x-www-form-urlencoded")

	if e != nil {
		t.Error(e)
	}
}

func TestDecodeBodyJSON(t *testing.T) {
	m1 := GT1{"Mercedes-Benz", "CLK GTR", 1997}
	b, e := json.Marshal(m1)

	if e != nil {
		t.Fatal(e)
	}

	reader := bytes.NewReader(b)
	r, e := http.NewRequest(http.MethodPost, "http://example.com", reader)

	if e != nil {
		t.Fatal(e)
	}

	r.Header.Set("Content-Type", "application/json")

	m2 := GT1{}
	e = DecodeBodyJSON(r, &m2)

	if e != nil {
		t.Fatal(e)
	}

	if m1.Manufacturer != m2.Manufacturer ||
		m1.Model != m2.Model ||
		m1.Debut != m2.Debut {

		t.Fatalf("%+v does not match decoded %+v", m1, m2)
	}

	r.Header.Set("Content-Type", "text/plain")

	m2 = GT1{}
	e = DecodeBodyJSON(r, &m2)

	if e == nil {
		t.Fatal("Decoded body with incorrect content type")
	}
}

func checkContentType(t *testing.T, res *http.Response) {
	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("SendJSON set incorrect Content-Type header: %s", ct)
	}
}
