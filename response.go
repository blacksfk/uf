package microframework

import (
	"net/http"
	"encoding/json"
)

type ResponseWriter struct {
	http.ResponseWriter
}

// send a JSON response (200 OK)
func (w ResponseWriter) JSON(data interface{}) error {
	encoder := json.NewEncoder(w)

	// set the headers
	w.Header().Set("Content-Type", "application/json")

	// attempt to encode the data as JSON
	e := encoder.Encode(data)

	if e != nil {
		// invalid data, or some other error
		return e
	}

	return nil
}

// send a JSON error response (he.Code)
func (w ResponseWriter) ErrorJSON(he HttpError) error {
	encoder := json.NewEncoder(w)

	// write the headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(he.Code)

	// attempt to encode the http error as JSON
	e := encoder.Encode(he)

	if e != nil {
		// something went incredibly wrong
		return e
	}

	return nil
}
