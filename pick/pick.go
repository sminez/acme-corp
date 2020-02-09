// pick - a minimalist input selector
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"9fans.net/go/acme"
	"9fans.net/go/draw"
)

const (
	windowTitle   = "pick"
	defaultPrompt = "> "
)

var (
	ignoreCase = flag.Bool("i", false, "ignore case")
	allowRegex = flag.Bool("r", false, "allow regular expression input")
	promtStr   = flag.String("p", defaultPrompt, "prompt to display before input")
)

// TODO: needs to actually compute dimensions for the window rather than return a default
//   - will need to be based on selected font size and lineCount
func getWindowSize() string {
	return "1024x800"
}

// runs in a goroutine
func handleKeyboardInput(d *draw.Display, ch chan string) {
	var kbCtl *draw.Keyboardctl
	var currentInput string

	kbCtl = d.InitKeyboard()
	for r := range kbCtl.C {
		if r <= 26 {
			switch fmt.Sprintf("C-%c", r+96) {
			case "C-h":
				l := len(currentInput)
				if l > 0 {
					currentInput = currentInput[:l-1]
				}

			case "C-w":
				words := strings.Split(currentInput, " ")
				currentInput = strings.Join(words[:len(words)-1], " ")

			case "C-c", "C-d":
				os.Exit(0)

			default:
				// unknown control sequence so drop it
			}
		} else {
			currentInput += string(r)
		}

		ch <- currentInput
	}
}

func initDisplay(font, windowSize string) *draw.Display {
	var display *draw.Display
	var err error

	if display, err = draw.Init(nil, font, windowTitle, windowSize); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR unable to open display: %s\n", err)
		os.Exit(1)
	}
	return display
}

// TODO: should this be a goro that takes a channel of filter strings
// and outputs to a channel of filtered lines? The main loop can then
// send a new filter to reset and get new lines provided it clears the
// input first.
func filterLines(filter string, lines []string) []string {

}

func main() {
	var lines []string
	var win *acme.Win
	var err error
	flag.Parse()

	if win, err = acme.New(); err != nil {
		fmt.Printf("Unable to initialise new acme window: %s\n", err)
		os.Exit(1)
	}

	win.Name("+pick")
	win.Write("tag", []byte("Reset"))

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		lines = append(lines, s.Text())
	}

	fmt.Println("got: ", lines)

	disp := initDisplay(*fontSpec, getWindowSize())
	ch := make(chan string)

	go handleKeyboardInput(disp, ch)
	disp.Flush()

	// probably want to be smart about responding to edits in the filter
	// but for now we simply re-filter the input on every edit
	for filter := range ch {
		fmt.Printf("current filter: '%s'\n", filter)
		for _, l := range filterLines(filter, lines) {
			fmt.Println(l)
		}
	}
}
