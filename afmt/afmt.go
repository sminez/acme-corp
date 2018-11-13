// Package afmt watches acme for know file extensions on files
// being written inside acme. Each time a known file is written,
// it runs the appropriate Tool and reloads the file in acme.
//
// NOTE: This will _not_ run for files without the correct file extension
//       For example, python scripts without a .py
// TODO: Rewrite this to modify the _window_ body rather than the underlying
//		 files. Would this also require a check that we had been idempotent?
package afmt

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"

	"9fans.net/go/acme"
)

// A Tool is a program that can rewrite source files or report
// on errors that were encountered in the code.
type Tool struct {
	cmd          string
	args         []string
	outputFixer  func(string) string
	ignoreOutput bool
}

// TODO: copy the window body to a temp file, format it there and then
//       reload it in the window body. Probably best to set up the temp
//       file in FileType.Reformat and then pass that in here?
//		 >> Look at how acmego does this
func (t *Tool) reformat(e *acme.LogEvent) {
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

	if t.ignoreOutput {
		return
	}

	if len(output) > 0 {
		msg = fmt.Sprint(string(output))
		if t.outputFixer != nil {
			msg = t.outputFixer(msg)
		}
		fmt.Fprintf(os.Stdout, "%s: %s\n", t.cmd, msg)
	}
}

// A FileType defines a set of Tools and an associated file type
// to run them on. If the files support a unix shebang then we
// try to parse that as well if the extension is missing.
type FileType struct {
	extensions   []string
	shebangProgs []string
	Tools        []Tool
}

// Matches checks to see if this is a file we need to reformat
// TODO: parse shebangs!
func (f *FileType) Matches(e *acme.LogEvent) bool {
	fileExtension := path.Ext(e.Name)
	for _, ext := range f.extensions {
		// remove the .
		if ext == fileExtension[1:] {
			return true
		}
		// if shebang is correct, return true
	}
	return false
}

// Reformat applies all known formatters to the underlying file
func (f *FileType) Reformat(e *acme.LogEvent) {
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

// WatchAndFormat is an example event loop using afmt
func WatchAndFormat() {
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
				if ft.Matches(&event) {
					ft.Reformat(&event)
				}
			}
		}
	}
}
