package microframework

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type GT1 struct {
	manufacturer, model string
	debut               int
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
	b, e := ioutil.ReadAll(res.Body)

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
	b, e := ioutil.ReadAll(res.Body)

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
	type testCase struct {
		method, contentType string
		body                io.Reader
	}

	b, e := json.Marshal(GT1{"Toyota", "TS020", 1998})

	if e != nil {
		t.Fatal(e)
	}

	reader := bytes.NewReader(b)
	passCases := []testCase{
		// test methods which have request bodies
		{http.MethodPost, "application/json", reader},
		{http.MethodPut, "application/json", reader},
		{http.MethodPatch, "application/json", reader},
	}

	// loop through and ensure no error has occurred
	for _, c := range passCases {
		r, e := http.NewRequest(c.method, "http://example.com", c.body)

		if e != nil {
			t.Fatal(e)
		}

		r.Header.Set("Content-Type", c.contentType)

		_, e = ReadBody(r)

		if e != nil {
			t.Error(e)
		}
	}

	reader = bytes.NewReader([]byte(`<<<<>>>>`))
	failCases := []testCase{
		// test methods which have request bodies
		{http.MethodPost, "application/pdf", reader},
		{http.MethodPut, "text/text", reader},
		{http.MethodPatch, "garbage", reader},
	}

	// loop through and ensure an error has occurred
	for _, c := range failCases {
		r, e := http.NewRequest(c.method, "http://example.com", c.body)

		if e != nil {
			t.Fatal(e)
		}

		r.Header.Set("Content-Type", c.contentType)

		j, e := ReadBody(r)

		if e == nil {
			t.Errorf("Invalid JSON decoded: %s", string(j))
		}
	}
}

func checkContentType(t *testing.T, res *http.Response) {
	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("SendJSON set incorrect Content-Type header: %s", ct)
	}
}
