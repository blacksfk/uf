package microframework

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
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
// header does not match any of the provided content types.
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

// Decode the request body into ptr. Returns a 400 Bad Request error if the
// received Content-Type header is not application/json.
func DecodeBodyJSON(r *http.Request, ptr interface{}) error {
	bytes, e := ReadBody(r, "application/json")

	if e != nil {
		return e
	}

	return json.Unmarshal(bytes, ptr)
}

// Get a URL parameter.
//
// In order to test handlers that require parameters to operate,
// a mock request should be created and the context manipulated so that
// the params are embedded within the context.
// Example (error handling omitted for the sake of brevity):
//
// // URL: example.com/mk/:char/result/:tournament
// params := httprouter.Params{
// 		{Key: "char", Value: "Kitana"},
// 		{Key: "tournament", Value: "10"},
// }
//
// // create and embed the context in the request
// ctx := context.WithValue(context.Background(), httprouter.ParamsKey, params)
// r, _ := http.NewRequestWithContext(ctx, http.MethodGet, "example.com", nil)
// w := http.NewRecorder()
//
// controller.Handle(w, r)
//
// // determine whether the test case passes...
func GetParam(r *http.Request, name string) string {
	params := httprouter.ParamsFromContext(r.Context())

	return params.ByName(name)
}

// Get a URL parameter as an int64. See GetParam for an example on how
// to test parameters in handlers.
func GetParamInt64(r *http.Request, name string) (int64, error) {
	str := GetParam(r, name)

	return strconv.ParseInt(str, 10, 0)
}

// Get a URL parameter type-cast as an int. See GetParam for an example
// on how to test parameters in handlers.
func GetParamInt(r *http.Request, name string) (int, error) {
	i64, e := GetParamInt64(r, name)

	return int(i64), e
}
