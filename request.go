package microframework

import (
	"io/ioutil"
	"net/http"
	"encoding/json"
	"strconv"
)

type Request struct {
	*http.Request
	JSONBody interface{}
}

// get a URL query parameter
func (r *Request) GetParam(name string) string {
	return r.URL.Query().Get(":" + name)
}

// get a URL query parameter type-cast as an int
func (r *Request) GetParamInt(name string) (int, error) {
	paramString := r.URL.Query().Get(":" + name)
	paramInt64, e := strconv.ParseInt(paramString, 10, 0)

	if e != nil {
		return 0, e
	}

	return int(paramInt64), nil
}

// parse the body of the request as JSON. Returns an error if the received
// Content-Type is not "application/json".
func (r *Request) parseBody() error {
	bytes, e := ioutil.ReadAll(r.Body)

	if e != nil {
		return e
	}

	length := len(bytes)
	ct := r.Header.Get("Content-Type")

	if length == 0 {
		// no body to unmarshal
		return nil
	} else if length > 0 && ct != "application/json" {
		// only unmarshal if body is present and json
		return BadRequest("Bad Content-Type: " + ct)
	}

	var blob interface{}
	e = json.Unmarshal(bytes, &blob)

	if e != nil {
		return e
	}

	r.JSONBody = blob

	return nil
}
