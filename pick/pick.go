// pick - a minimalist input selector
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"9fans.net/go/acme"
)

const promptStr = "> "

var (
	ignoreCase          = flag.Bool("i", false, "ignore case")
	snarfCurrentWindow  = flag.Bool("a", false, "use the current acme window as input and jump to the selection")
	getWinidFromSnooper = flag.Bool("s", false, "get the current windowID from the snooper not the environment")
)

type linePicker struct {
	w             *acme.Win
	lines         []string
	filteredLines []string
	currentInput  string
	ignoreCase    bool
}

func newLinePicker(lines []string, ignoreCase bool) *linePicker {
	w, err := acme.New()
	if err != nil {
		fmt.Printf("Unable to initialise new acme window: %s\n", err)
		os.Exit(1)
	}
	w.Name("+pick")
	w.Write("tag", []byte("Reset"))

	return &linePicker{
		w:             w,
		lines:         lines,
		filteredLines: lines,
		currentInput:  "",
		ignoreCase:    ignoreCase,
	}
}

func (lp *linePicker) filter() (int, string, error) {
	var result string

	for e := range lp.w.EventChan() {
		reRender := false

		if e.C1 != 'K' {
			// TODO: need to allow actions in the tag and maybe selecting the a line
			//       via button 3, but probably should drop button 2 events?
			lp.w.WriteEvent(e)
			continue
		}

		fmt.Printf("%#v\n", e)

		r := e.Text[0]
		if r <= 26 {
			switch fmt.Sprintf("C-%c", r+96) {
			case "C-h": // backspace
				l := len(lp.currentInput)
				if l > 0 {
					lp.currentInput = lp.currentInput[:l-1]
				}
				reRender = true

			case "C-j": // return
				break

			case "C-w": // backwards kill word
				words := strings.Split(lp.currentInput, " ")
				lp.currentInput = strings.Join(words[:len(words)-1], " ")
				reRender = true

			case "C-c", "C-d":
				lp.w.Del(true)
				os.Exit(0)

			default:
				// unknown control sequence so drop it
			}
		} else {
			lp.currentInput += string(r)
			reRender = true

		}

		if reRender {
			// TODO: now trim lines to only show what has been filtered
			// under the input line (that needs init-ing as well)
			msg := fmt.Sprintf("input: %s\n", lp.currentInput)
			lp.w.Write("errors", []byte(msg))
		}
	}

	return 0, result, nil
}

func getCurrentWindow() (*acme.Win, error) {
	winStr := os.Getenv("winid")
	if len(winStr) == 0 {
		return nil, fmt.Errorf("-a can only be used from inside acme")
	}

	winID, err := strconv.Atoi(winStr)
	if err != nil {
		return nil, fmt.Errorf("non numeric winid: %s", winStr)
	}

	w, err := acme.Open(winID, nil)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func main() {
	var (
		lines []string
		body  []byte
		w     *acme.Win
		err   error
	)

	flag.Parse()

	if *snarfCurrentWindow {
		if w, err = getCurrentWindow(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer w.CloseFiles()

		w.Addr(",")
		if body, err = w.ReadAll("data"); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		lines = strings.Split(string(body), "\n")

	} else {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			lines = append(lines, s.Text())
		}
	}

	lp := newLinePicker(lines, *ignoreCase)
	n, selection, err := lp.filter()
	fmt.Println(n, selection)
}
