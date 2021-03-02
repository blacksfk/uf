package microframework

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// Send a JSON response
func SendJSON(w http.ResponseWriter, data interface{}) error {
	encoder := json.NewEncoder(w)

	// set the headers
	w.Header().Set("Content-Type", "application/json")

	// attempt to encode the data as JSON
	return encoder.Encode(data)
}

// Send an HttpError as a JSON response
func SendErrorJSON(w http.ResponseWriter, he HttpError) error {
	encoder := json.NewEncoder(w)

	// write the headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(he.Code)

	// attempt to encode the http error as JSON
	return encoder.Encode(he)
}

// Returns the bytes read from r.Body. Returns a Bad Request error if the received Content-Type
// header is does not match any of the provided content types.
func ReadBody(r *http.Request, contentTypes ...string) ([]byte, error) {
	if l := len(contentTypes); l > 0 {
		valid := false
		ct := r.Header.Get("Content-Type")

		for i := 0; i < l; i++ {
			if ct == contentTypes[i] {
				valid = true
				break
			}
		}

		if !valid {
			b := strings.Builder{}
			b.WriteString("Bad Content-Type: ")
			b.WriteString(ct)
			b.WriteString(". Accept: ")
			b.WriteString(strings.Join(contentTypes, ", "))

			return nil, BadRequest(b.String())
		}
	}

	defer r.Body.Close()
	return io.ReadAll(r.Body)
}

func DecodeBodyJSON(r *http.Request, ptr interface{}) error {
	bytes, e := ReadBody(r, "application/json")

	if e != nil {
		return e
	}

	return json.Unmarshal(bytes, ptr)
}

// Get a URL query parameter
func GetParam(r *http.Request, name string) string {
	return r.URL.Query().Get(":" + name)
}

// Get a URL query parameter as an int64
func GetParamInt64(r *http.Request, name string) (int64, error) {
	str := r.URL.Query().Get(":" + name)

	return strconv.ParseInt(str, 10, 0)
}

// Get a URL query parameter type-cast as an int
func GetParamInt(r *http.Request, name string) (int, error) {
	i64, e := GetParamInt64(r, name)

	return int(i64), e
}
