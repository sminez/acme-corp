package acorp

import (
	"fmt"
	"strings"

	"9fans.net/go/acme"
)

// SetCursorEOL will position the current window cursor at the end of line.
func SetCursorEOL(w *acme.Win, line int) {
	w.Addr(fmt.Sprintf("%d-#1", line+1))
	w.Ctl("dot=addr")
	w.Ctl("show")
}

// SetCursorBOL will position the current window cursor at the beginning of line.
func SetCursorBOL(w *acme.Win, line int) {
	w.Addr(fmt.Sprintf("%d-#0", line))
	w.Ctl("dot=addr")
	w.Ctl("show")
}

// WindowBody reads the body of the current window as single string
func WindowBody(w *acme.Win) (string, error) {
	var (
		body []byte
		err  error
	)

	// TODO: stash and restore current addr
	w.Addr(",")
	if body, err = w.ReadAll("data"); err != nil {
		return "", err
	}
	return string(body), nil
}

// WindowBodyLines reads the body of the current window as an array of strings split on newline
func WindowBodyLines(w *acme.Win) ([]string, error) {
	body, err := WindowBody(w)
	if err != nil {
		return nil, err
	}
	return strings.Split(body, "\n"), nil
}

// EventLineNumber returns the line that a e occurred on in w
func EventLineNumber(w *acme.Win, e *acme.Event) (int, error) {
	body, err := WindowBody(w)
	if err != nil {
		return -1, err
	}

	upToCursor := body[:e.Q0]
	return strings.Count(upToCursor, "\n"), nil
}
