package microframework

import (
	"fmt"
	"net/http"
	"time"
)

type queue struct {
	c  Handler
	m  []Middleware
	el ErrorLogger
	al AccessLogger
}

// Create a new queue.
func newQueue(c Handler, m []Middleware, el ErrorLogger, al AccessLogger) *queue {
	return &queue{c, m, el, al}
}

// Queue implements http.Handler.
func (q *queue) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
func (q *queue) logAccess(r *http.Request, start time.Time) {
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

func (q *queue) handleError(w http.ResponseWriter, e error) {
	// check if e is already an http error
	httpError, ok := e.(HttpError)

	if !ok {
		// create an HttpError out of the plain error
		httpError = InternalServerError(e.Error())
	}

	// log the error used the supplied function
	q.el(httpError)

	// send the error to the client
	e = SendErrorJSON(w, httpError)

	if e != nil {
		// something went incredibly wrong...
		q.el(fmt.Errorf("SendErrorJSON(): %v", e))
	}
}
