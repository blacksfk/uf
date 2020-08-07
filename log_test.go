package microframework

import (
	"bytes"
	"io"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	r, w, e := os.Pipe()

	if e != nil {
		t.Fatalf("os.Pipe(): %v", e)
	}

	req := httptest.NewRequest(http.MethodGet, "http://www.example.com/", nil)

	// verify that the request method and uri were printed, with any duration
	re, e := regexp.Compile(fmt.Sprintf(`%s\s%s\s\d+[mun]s`, req.Method, req.RequestURI))

	if e != nil {
		w.Close()
		t.Fatalf("regexp.Compile(): %v", e)
	}

	// maintain a copy of stdout and set the pipe to
	// the write pipe created above
	og := os.Stdout

	// capture the output
	os.Stdout = w
	logReq(req, time.Now())

	// copy bytes from the pipe to something readable using a go routine
	out := make(chan string)
	err := make(chan error)

	go copyBuffer(r, out, err)

	// close the pipes and reset os.Stdout
	e = w.Close()

	if e != nil {
		t.Fatalf("w.Close(): %v", e)
	}

	os.Stdout = og

	// wait for response on one of the channels
	// if waiting for too long, fail the test
	attempts := 0
	max := 5
	waitInc := time.Millisecond * 100

	select {
	case str := <-out:
		if !re.MatchString(str) {
			t.Errorf("Expected (regular expression): %s, actual: %s\n", re, str)
		}
	case e = <-err:
		t.Fatalf("io.Copy(): %v", e)
	default:
		attempts++

		if attempts > max {
			total := waitInc.Milliseconds() * int64(attempts)

			t.Fatalf("Did not receive a message on any channel in %dms", total)
		}

		time.Sleep(waitInc)
	}
}

func copyBuffer(r *os.File, out chan string, err chan error) {
	buf := bytes.NewBuffer([]byte{})
	_, e := io.Copy(buf, r)

	if e != nil {
		err <- e
	} else {
		out <- buf.String()
	}
}
