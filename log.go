package microframework

import (
	"fmt"
	"net/http"
	"time"
)

// Logs requests to stdout. Format: "METHOD URI DURATION"
func LogRequest(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		start := time.Now()

		defer logReq(r, start)
		return next(w, r)
	}
}

func logReq(r *http.Request, start time.Time) {
	var duration int64
	var unit string
	magnitude := time.Since(start)

	// use d != 0 in case start is in the future making magnitude negative
	if d := magnitude.Milliseconds(); d != 0 {
		duration = d
		unit = "m"
	} else if d := magnitude.Microseconds(); d != 0 {
		duration = d
		unit = "u"
	} else {
		duration = magnitude.Nanoseconds()
		unit = "n"
	}

	fmt.Printf("%s %s %d%ss\n", r.Method, r.RequestURI, duration, unit)
}
