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

	"9fans.net/go/acme"
)

// FTYPES lists the currently known filetypes for afmt
var FTYPES = []FileType{python, golang, rust, shell, javascript, json}

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
	fileExtension := path.Ext(e.Name)
	for _, ext := range f.Extensions {
		if ext == fileExtension {
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

var python = FileType{
	Extensions:   []string{"py", "pyw"},
	ShebangProgs: []string{"python"},
	Tools: []Tool{
		Tool{Cmd: "isort", Args: []string{"-m", "5"}},
		Tool{Cmd: "black"},
		Tool{Cmd: "flake8"},
	},
}

var golang = FileType{
	Extensions: []string{"go"},
	Tools: []Tool{
		Tool{Cmd: "gofmt", Args: []string{"-w"}},
		Tool{Cmd: "goimports", Args: []string{"-w"}},
		Tool{Cmd: "go", Args: []string{"vet"}},
	},
}

var rust = FileType{
	Extensions: []string{"rs"},
	Tools:      []Tool{Tool{Cmd: "rustfmt"}},
}

var shell = FileType{
	Extensions: []string{"sh", "bash", "zsh"},
	Tools: []Tool{
		// Remove trailing whitespace and whitespace only lines
		Tool{Cmd: "sed", Args: []string{"-i", "'s/[[:blank:]]*$//g'"}},
		Tool{Cmd: "shellcheck", Args: []string{"--color=never"}},
	},
}

var javascript = FileType{
	Extensions: []string{"js"},
	Tools: []Tool{
		Tool{Cmd: "js-beautify", Args: []string{"-r"}},
		Tool{Cmd: "jshint"},
	},
}

var json = FileType{
	Extensions: []string{"json"},
	Tools:      []Tool{Tool{Cmd: "json-format"}},
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
