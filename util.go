package microframework

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

// send a JSON response
func SendJSON(w http.ResponseWriter, data interface{}) error {
	encoder := json.NewEncoder(w)

	// set the headers
	w.Header().Set("Content-Type", "application/json")

	// attempt to encode the data as JSON
	return encoder.Encode(data)
}

// send an HttpError as a JSON response
func SendErrorJSON(w http.ResponseWriter, he HttpError) error {
	encoder := json.NewEncoder(w)

	// write the headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(he.Code)

	// attempt to encode the http error as JSON
	return encoder.Encode(he)
}

// returns the bytes read from r.Body. Returns an error if the received Content-Type
// header is not "application/json".
func ReadJSON(r *http.Request) ([]byte, error) {
	if ct := r.Header.Get("Content-Type"); ct != "application/json" {
		return nil, BadRequest("Bad Content-Type: " + ct)
	}

	return ioutil.ReadAll(r.Body)
}

// get a URL query parameter
func GetParam(r *http.Request, name string) string {
	return r.URL.Query().Get(":" + name)
}

// get a URL query parameter type-cast as an int
func GetParamInt(r *http.Request, name string) (int, error) {
	str := r.URL.Query().Get(":" + name)
	i64, e := strconv.ParseInt(str, 10, 0)

	if e != nil {
		return 0, e
	}

	return int(i64), nil
}
