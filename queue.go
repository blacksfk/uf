package microframework

import (
	"fmt"
	"net/http"
)

type queue struct {
	c  Handler
	m  []Middleware
	el ErrorLogger
}

// Create a new queue.
func newQueue(c Handler, m []Middleware, el ErrorLogger) *queue {
	return &queue{c, m, el}
}

// Queue implements http.Handler.
func (q *queue) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var e error = nil

	for _, m := range q.m {
		if e = m(r); e != nil {
			q.handleError(w, e)

			return
		}
	}

	if e = q.c(w, r); e != nil {
		q.handleError(w, e)
	}
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
