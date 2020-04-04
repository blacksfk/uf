package microframework

import (
	"fmt"
	"net/http"
)

type queue struct {
	first Handler
}

// create a new queue
func newQueue(controller Handler, middleware []Middleware) *queue {
	var next Handler
	curr := controller

	// call the middleware function with the result of the next middleware
	// handler as a parameter. Starts from the end, going in reverse order
	// with the controller as the parameter to the last middleware
	for i := len(middleware)-1; i >= 0; i-- {
		next = curr
		curr = middleware[i](next)
	}

	return &queue{curr}
}


// make queue implement http.Handler
func (q *queue) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// wrap the response with extra methods
	customResponse := ResponseWriter{w}

	// wrap the request and parse the body
	customRequest := Request{r, nil}
	e := customRequest.parseBody()

	if e != nil {
		handleError(&customResponse, e)

		return
	}

	// start calling functions in the queue
	e = q.first(&customResponse, &customRequest)

	if e != nil {
		handleError(&customResponse, e)
	}
}

func handleError(crw *ResponseWriter, e error) {
	fmt.Println(e)

	// check if e is already an http error
	httpError, ok := e.(HttpError)

	if !ok {
		// create an HttpError out of the plain error
		httpError = InternalServerError(e.Error())
	}

	// send the error to the client
	e = crw.ErrorJSON(httpError)

	if e != nil {
		// something went incredibly wrong...
		fmt.Println(e)
	}
}
