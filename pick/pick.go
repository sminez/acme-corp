package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"9fans.net/go/draw"
)

const (
	windowTitle       = "pick"
	defaultFont       = "/mnt/font/TerminessTTFNerdFontComplete-Medium/12a/font"
	defaultPrompt     = "> "
	defaultLineCount  = 10
	defaultNormalBg   = "#f2e5bc"
	defaultSelectedBg = "#a89984"
	defaultNormalFg   = "#282828"
	defaultSelectedFg = "#eeeeee"
)

var (
	ignoreCase = flag.Bool("i", false, "ignore case")
	allowRegex = flag.Bool("r", false, "allow regular expression input")
	lineCount  = flag.Int("l", defaultLineCount, "lines to display at once")
	promtStr   = flag.String("p", defaultPrompt, "prompt to display before input")
	fontSpec   = flag.String("fn", defaultFont, "font name")
	normalBg   = flag.String("nb", defaultNormalBg, "normal background color in #RRGGBB format")
	selectedBg = flag.String("sb", defaultSelectedBg, "selected background color in #RRGGBB format")
	normalFg   = flag.String("nf", defaultNormalFg, "normal foreground color in #RRGGBB format")
	selectedFg = flag.String("sf", defaultSelectedFg, "selected foreground color in #RRGGBB format")
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
		}

		currentInput += string(r)
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

func main() {
	flag.Parse()

	// var lines []string
	// s := bufio.NewScanner(os.Stdin)
	// for s.Scan() {
	// 	lines = append(lines, s.Text())
	// }

	disp := initDisplay(*fontSpec, getWindowSize())
	ch := make(chan string)

	go handleKeyboardInput(disp, ch)
	disp.Flush()

	for input := range ch {
		fmt.Printf("current input: '%s'\n", input)
	}
}
