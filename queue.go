package microframework

import (
	"fmt"
	"net/http"
)

type queue struct {
	first Handler
	el ErrorLogger
}

// Create a new queue
func newQueue(controller Handler, middleware []Middleware, el ErrorLogger) *queue {
	curr := controller

	// call the middleware function with the result of the next middleware
	// handler as a parameter. Starts from the end, going in reverse order
	// with the controller as the parameter to the last middleware
	for i := len(middleware) - 1; i >= 0; i-- {
		next := curr
		curr = middleware[i](next)
	}

	return &queue{curr, el}
}

// Queue implements http.Handler
func (q *queue) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// start calling functions in the queue
	e := q.first(w, r)

	if e != nil {
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
