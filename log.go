package microframework

import (
	"fmt"
	"net/http"
)

// Logs requests to stdout. Format: "method uri duration(unit)s"
func LogStdout(r *http.Request, duration int64, unit string) {
	fmt.Printf("%s %s %d%ss\n", r.Method, r.RequestURI, duration, unit)
}
