// acmefmt watches acme for know file extensions on files
// being written inside acme. Each time a known file is written,
// it runs the appropriate tool and reloads the file in acme.
//
// NOTE: This will _not_ run for files without the correct file extension
//       For example, python scripts without a .py
// TODO: Parse file shebangs to determine filetype as a fallback.
//		 Look at https://godoc.org/golang.org/x/tools for more go tools
//	       - go fix, go guru
//		 Add: Rust, HTML, CSS, JS, Java, Kotlin, Shell
// TODO: Rewrite this to modify the _window_ body rather than the underlying
//		 files. Would this also require a check that we had been idempotent?
// TODO: Make URLs part of the tool struct and use this to install them? (maybe not...)
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"

	"9fans.net/go/acme"
)

// A Tool is a program that can rewrite source files or report
// on errors that were encountered in the code.
type Tool struct {
	Cmd            string
	Args           []string
	OutputFixer    func(string) string
	RewritesSource bool
}

// TODO: copy the window body to a temp file, format it there and then
//       reload it in the window body. Probably best to set up the temp
//       file in FileType.reformat and then pass that in here?
//		 >> Look at how acmego does this
func (t *Tool) reformat(e *acme.LogEvent) {
	_, err := ioutil.ReadFile(e.Name)
	if err != nil {
		return
	}

	args := append(t.Args, e.Name)
	output, err := exec.Command(t.Cmd, args...).CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s error: %s\n", t.Cmd, err)
	}
	if output != nil {
		fmt.Fprintf(os.Stdout, "%s: %s\n", t.Cmd, output)
	}
}

// A FileType defines a set of tools and an associated file type
// to run them on. If the files support a unix shebang then we
// try to parse that as well if the extension is missing.
type FileType struct {
	Extensions   []string
	ShebangProgs []string
	Tools        []Tool
}

// Check to see if this is a file we need to reformat
// TODO: parse shebangs!
func (f *FileType) matches(e *acme.LogEvent) bool {
	for _, ext := range f.Extensions {
		if strings.HasSuffix(e.Name, ext) {
			return true
		}
		// if shebang is correct, return true
	}
	return false
}

// Apply all formatters to the underlying file
func (f *FileType) reformat(e *acme.LogEvent) {
	w, err := acme.Open(e.ID, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer w.CloseFiles()

	for _, t := range f.Tools {
		t.reformat(e)
	}
	w.Write("ctl", []byte("get"))
}

// Python external tooling for acme.
// automatic import ordering and grouping with isort, overall formatting
// with black and then syntax and style linting with flake8
var Python = FileType{
	Extensions:   []string{"py", "pyw"},
	ShebangProgs: []string{"python"},
	Tools: []Tool{
		Tool{
			Cmd:  "isort", // https://github.com/timothycrosley/isort
			Args: []string{"-m", "5"},
		},
		Tool{
			Cmd:  "black", // https://github.com/ambv/black
			Args: []string{""},
		},
		Tool{
			Cmd:  "flake8", // https://gitlab.com/pycqa/flake8
			Args: []string{""},
		},
	},
}

// Golang external tooling for acme.
// Overall formatting with gofmt, import management with goimports
// and linting for common errors with govet.
var Golang = FileType{
	Extensions:   []string{"go"},
	ShebangProgs: []string{""},
	Tools: []Tool{
		Tool{
			Cmd:  "gofmt", // https://golang.org/cmd/gofmt/
			Args: []string{"-w"},
		},
		Tool{
			Cmd:  "goimports", // https://godoc.org/golang.org/x/tools/cmd/goimports
			Args: []string{"-w"},
		},
		Tool{
			Cmd:  "go", // https://godoc.org/golang.org/x/tools/cmd/govet
			Args: []string{"vet"},
		},
	},
}

// WatchAndFix is intended to be run in it's own goroutine or as a stand alone
// main function. It watches the acme event log for know file types being saved
// and then runs the provided external tools over the window content.
func WatchAndFix(fileTypes []FileType) {
	l, err := acme.Log()
	if err != nil {
		log.Fatal(err)
	}

	for {
		event, err := l.Read()
		if err != nil {
			log.Fatal(err)
		}

		// On window save, run any tools we know about
		if event.Name != "" && event.Op == "put" {
			for _, ft := range fileTypes {
				if ft.matches(&event) {
					ft.reformat(&event)
				}
			}
		}
	}
}

func main() {
	// TODO: add flag for whether or not we should format
	WatchAndFix([]FileType{Python, Golang})
}
