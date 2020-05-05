/*
pick - a minimalist input selector for the acme text editor modelled after dmenu

If launched from within acme itself, pick will use the current acme window as defined by the 'winid'
environment variable. Otherwise, it will attempt to query a running snooper instance to fetch the
focused window id.
  * To mimic dmenu behaviour of reading input from stdin, pass the '-s' flag.
  * To return the index of the selected line in the input instead of the line itself, pass the '-n' flag.
  * To override the default prompt ('> ') pass the '-p' flag followed by the string to use as the prompt.

+pick window actions
  * character input will be interpreted as a regex for filtering lines
  * button 3: select a line to return
  * Return:   select the line the curesor is currently on
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"9fans.net/go/acme"
	"github.com/sminez/acme-corp/acorp"
)

const (
	lineOffset = 1 // assuming single line prompt
	windowName = "+pick"
)

var (
	readFromStdIn = flag.Bool("s", false, "read input from stdin instead of the current acme window")
	returnLineNum = flag.Bool("n", false, "return the line number of the selected line, not the line itself")
	numberLines   = flag.Bool("N", false, "prefix each line with its line number")
	prompt        = flag.String("p", "> ", "prompt to present to the user when taking input")
)

type linePicker struct {
	w              *acme.Win
	rawLines       []string
	lineMap        map[int]string // original line numbers
	selectedLines  map[int]int    // window line numbers -> input line number
	currentInput   string
	selectionEvent *acme.Event
}

func newLinePicker(rawLines []string) *linePicker {
	var w *acme.Win
	var err error

	if w, err = acme.New(); err != nil {
		fmt.Printf("Unable to initialise new acme window: %s\n", err)
		os.Exit(1)
	}

	w.Name(windowName)
	lineMap := make(map[int]string)

	for n, l := range rawLines {
		lineMap[n+1] = l
	}

	return &linePicker{
		w:             w,
		rawLines:      rawLines,
		lineMap:       lineMap,
		selectedLines: make(map[int]int),
		currentInput:  "",
	}
}

func (lp *linePicker) selectedLine() (int, string, error) {
	var (
		windowLineNumber int
		err              error
	)

	if windowLineNumber, err = acorp.EventLineNumber(lp.w, lp.selectionEvent); err != nil {
		return -1, "", err
	}

	// hitting enter on the input line selects top match if there is at least one, otherwise
	// it returns the current input text
	if windowLineNumber == 0 {
		if len(lp.selectedLines) == 0 {
			return -1, lp.currentInput, nil
		}
		windowLineNumber = 1
	}

	lineNumber := lp.selectedLines[windowLineNumber-lineOffset]
	return lineNumber, lp.lineMap[lineNumber], nil
}

func (lp *linePicker) filter() (int, string, error) {
	ef := &acorp.EventFilter{
		KeyboardInputBody: func(w *acme.Win, e *acme.Event, done func() error) error {
			r := e.Text[0]

			if r <= 26 {
				switch fmt.Sprintf("C-%c", r+96) {
				case "C-j": // (Enter)
					lp.selectionEvent = e
					return done()
				case "C-d":
					lp.w.Del(true)
					os.Exit(0)
				}
			}

			lp.currentInput += string(r)
			return lp.reRender()
		},

		KeyboardDeleteBody: func(w *acme.Win, e *acme.Event, done func() error) error {
			if l := len(lp.currentInput); l > 0 {
				removed := e.Q1 - e.Q0
				lp.currentInput = lp.currentInput[:l-removed]
			}

			return lp.reRender()
		},

		Mouse3Body: func(w *acme.Win, e *acme.Event, done func() error) error {
			lp.selectionEvent = e
			return done()
		},
	}

	if err := lp.reRender(); err != nil {
		return -1, "", err
	}

	if err := ef.Filter(lp.w); err != nil {
		return -1, "", err
	}

	return lp.selectedLine()
}

func (lp *linePicker) reRender() error {
	lines := lp.rawLines
	lp.w.Clear()

	if len(lp.currentInput) > 0 {
		fragments := strings.Split(lp.currentInput, " ")
		lp.selectedLines = make(map[int]int)
		lines = []string{}
		k := 0

		for ix, line := range lp.lineMap {
			if containsAll(line, fragments) {
				lp.selectedLines[k] = ix
				lines = append(lines, line)
				k++
			}
		}
	}

	lp.w.Write("body", []byte(fmt.Sprintf("%s%s\n", *prompt, lp.currentInput)))
	lp.w.Write("body", []byte(strings.Join(lines, "\n")))
	acorp.SetCursorEOL(lp.w, 1)
	return nil
}

func containsAll(s string, ts []string) bool {
	for _, t := range ts {
		if !strings.Contains(s, t) {
			return false
		}
	}
	return true
}

func readFromAcme() ([]string, error) {
	var (
		w   *acme.Win
		err error
	)

	if w, err = acorp.GetCurrentWindow(); err != nil {
		return nil, err
	}
	defer w.CloseFiles()

	return acorp.WindowBodyLines(w)
}

func numberedLines(lines []string) []string {
	for ix, line := range lines {
		lines[ix] = fmt.Sprintf("%3d | %s", ix+1, line)
	}

	return lines
}

func main() {
	var (
		lines     []string
		err       error
		n         int
		selection string
	)

	flag.Parse()

	if *readFromStdIn {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			lines = append(lines, s.Text())
		}
	} else {
		if lines, err = readFromAcme(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if *numberLines {
		lines = numberedLines(lines)
	}

	lp := newLinePicker(lines)
	defer lp.w.Del(true)

	if n, selection, err = lp.filter(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if *returnLineNum {
		fmt.Println(n)
	} else {
		fmt.Println(selection)
	}
}
