package uf

import (
	"fmt"
	"net/http"
	"time"
)

// Queue handles errors returned from middleware and handlers
// along with implementing http.Handler. A Queue should not be
// created directly, and is only exposed to serve testing purposes
// with the httptest library (or alternatives) using the NewHttpTestHandler
// factory method.
type Queue struct {
	c  Handler
	m  []Middleware
	el ErrorLogger
	al AccessLogger
}

// Create a new queue.
func newQueue(c Handler, m []Middleware, config *Config) *Queue {
	return &Queue{c, m, config.ErrorLogger, config.AccessLogger}
}

// Create a test queue in order to use and test the uf.Handler
// with functions that only accept http.Handler. Eg. httptest.NewServer.
func NewHttpTestHandler(h Handler) http.Handler {
	return &Queue{c: h}
}

// Queue implements http.Handler.
func (q *Queue) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if q.al != nil {
		// access logger was provided
		start := time.Now()

		defer q.logAccess(r, start)
	}

	// loop through the middleware provided
	// terminate early if an error was returned
	for _, m := range q.m {
		if e := m(r); e != nil {
			q.handleError(w, e)

			return
		}
	}

	// run the controller function
	if e := q.c(w, r); e != nil {
		q.handleError(w, e)
	}
}

// Calculate the difference between start and now as an absolute value
// with an appropriate unit (ms, us, ns).
func (q *Queue) logAccess(r *http.Request, start time.Time) {
	var duration int64
	var unit string
	magnitude := time.Since(start)

	// use d != 0 in case start is in the future making magnitude negative
	// keep going until a non-zero d value is found (or just use nano seconds)
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

	q.al(r, duration, unit)
}

func (q *Queue) handleError(w http.ResponseWriter, e error) {
	// check if e is already an http error
	httpError, ok := e.(HttpError)

	if !ok {
		// create an HttpError out of the plain error
		httpError = InternalServerError(e.Error())
	}

	// log the error to the supplied function
	if q.el != nil {
		q.el(httpError)
	}

	// send the error to the client
	e = SendErrorJSON(w, httpError)

	if e != nil {
		// something went incredibly wrong...
		q.el(fmt.Errorf("SendErrorJSON(): %v", e))
	}
}
