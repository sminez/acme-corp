// afmt watches acme for know file extensions on files
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
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"

	"9fans.net/go/acme"
)

// FTYPES lists the currently known filetypes for afmt
var FTYPES = []fileType{
	python, golang, rust, shell, javascript, json,
}

// A tool is a program that can rewrite source files or report
// on errors that were encountered in the code.
type tool struct {
	cmd         string
	args        []string
	outputFixer func(string) string
}

// TODO: copy the window body to a temp file, format it there and then
//       reload it in the window body. Probably best to set up the temp
//       file in fileType.reformat and then pass that in here?
//		 >> Look at how acmego does this
func (t *tool) reformat(e *acme.LogEvent) {
	var msg string

	_, err := ioutil.ReadFile(e.Name)
	if err != nil {
		return
	}

	args := append(t.args, e.Name)
	output, err := exec.Command(t.cmd, args...).CombinedOutput()
	if err != nil {
		msg = err.Error()
		if t.outputFixer != nil {
			msg = t.outputFixer(msg)
		}
		fmt.Fprintf(os.Stderr, "%s error: %s\n", t.cmd, msg)
	}

	if output != nil {
		if t.outputFixer != nil {
			msg = fmt.Sprint(output)
			msg = t.outputFixer(msg)
		}
		fmt.Fprintf(os.Stdout, "%s: %s\n", t.cmd, msg)
	}
}

// A fileType defines a set of tools and an associated file type
// to run them on. If the files support a unix shebang then we
// try to parse that as well if the extension is missing.
type fileType struct {
	extensions   []string
	shebangProgs []string
	tools        []tool
}

// Check to see if this is a file we need to reformat
// TODO: parse shebangs!
func (f *fileType) matches(e *acme.LogEvent) bool {
	fileExtension := path.Ext(e.Name)
	for _, ext := range f.extensions {
		if ext == fileExtension {
			return true
		}
		// if shebang is correct, return true
	}
	return false
}

// Apply all formatters to the underlying file
func (f *fileType) reformat(e *acme.LogEvent) {
	w, err := acme.Open(e.ID, nil)
	if err != nil {
		log.Print(err)
		return
	}
	defer w.CloseFiles()

	for _, t := range f.tools {
		t.reformat(e)
	}
	w.Write("ctl", []byte("get"))
}

var python = fileType{
	extensions:   []string{"py", "pyw"},
	shebangProgs: []string{"python"},
	tools: []tool{
		tool{cmd: "isort", args: []string{"-m", "5"}},
		tool{cmd: "black"},
		tool{cmd: "flake8"},
	},
}

var golang = fileType{
	extensions: []string{"go"},
	tools: []tool{
		tool{cmd: "gofmt", args: []string{"-w"}},
		tool{cmd: "goimports", args: []string{"-w"}},
		tool{cmd: "go", args: []string{"vet"}},
	},
}

var rust = fileType{
	extensions: []string{"rs"},
	tools:      []tool{tool{cmd: "rustfmt"}},
}

var shell = fileType{
	extensions: []string{"sh", "bash", "zsh"},
	tools: []tool{
		// Remove trailing whitespace and whitespace only lines
		tool{cmd: "sed", args: []string{"-i", "'s/[[:blank:]]*$//g'"}},
		tool{cmd: "shellcheck", args: []string{"--color=never"}},
	},
}

var javascript = fileType{
	extensions: []string{"js"},
	tools: []tool{
		tool{cmd: "js-beautify", args: []string{"-r"}},
		tool{
			cmd: "jshint",
			outputFixer: func(s string) string {
				// Convert to button3 friendly output
				s = strings.Replace(s, " line ", "", -1)
				// Remove the column number as it just adds line noise
				re := regexp.MustCompile(", col .*,")
				return re.ReplaceAllString(s, " ->")
			},
		},
	},
}

var json = fileType{
	extensions: []string{"json"},
	tools:      []tool{tool{cmd: "json-format"}},
}

func main() {
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
			for _, ft := range FTYPES {
				if ft.matches(&event) {
					ft.reformat(&event)
				}
			}
		}
	}
}
