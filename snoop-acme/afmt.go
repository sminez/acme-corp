package snoop

// afmt watches acme for know file extensions on files being written inside
// acme. Each time a known file is written, it runs the appropriate Tool and
// reloads the file in acme.
//
// TODO: Rewrite this to modify the _window_ body rather than the underlying
//		 files. Would this also require a check that we had been idempotent?

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"

	"9fans.net/go/acme"
)

// A Tool is a program that can rewrite source files or report on errors that
// were encountered in the code.
type Tool struct {
	cmd            string
	args           []string
	outputFixer    func(string) string
	appendFilePath bool
	ignoreOutput   bool
}

// TODO: copy the window body to a temp file, format it there and then reload it
//       in the window body. Probably best to set up the temp file in
//       FileType.Reformat and then pass that in here? >> Look at how acmego
//       does this
func (t *Tool) reformat(e *acme.LogEvent) string {
	var args []string

	if _, err := ioutil.ReadFile(e.Name); err != nil {
		return err.Error()
	}

	if t.appendFilePath {
		args = append(t.args, e.Name)
	} else {
		args = append(t.args, path.Dir(e.Name))
	}

	b, _ := exec.Command(t.cmd, args...).CombinedOutput()
	if t.ignoreOutput || len(b) == 0 {
		return ""
	}

	if t.outputFixer != nil {
		return t.outputFixer(string(b))
	}

	return string(b)
}

// A FileType defines a set of Tools and an associated file type to run them on.
// If the files support a unix shebang then we try to parse that as well if the
// extension is missing.
type FileType struct {
	extensions   []string
	shebangProgs []string
	Tools        []Tool
}

// Matches checks to see if this is a file we need to reformat
func (f *FileType) Matches(e *acme.LogEvent) bool {
	fileExtension := path.Ext(e.Name)
	if len(fileExtension) > 0 {
		fileExtension = fileExtension[1:]
		for _, ext := range f.extensions {
			if fileExtension == ext {
				return true
			}
		}
	}

	s, err := getFirstLine(e.ID)
	if err != nil {
		return false
	}

	for _, prog := range f.shebangProgs {
		if strings.HasSuffix(s, prog) {
			return true
		}
	}

	return false
}

// Reformat applies all known formatters to the underlying file
func (f *FileType) Reformat(e *acme.LogEvent) string {
	var output string
	var w *acme.Win
	var err error

	if w, err = acme.Open(e.ID, nil); err != nil {
		fmt.Println(err)
		return err.Error()
	}
	defer w.CloseFiles()

	for _, t := range f.Tools {
		output += t.reformat(e)
	}

	w.Ctl("get")
	w.Ctl("clean")
	return output
}

func getFirstLine(winid int) (string, error) {
	w, err := acme.Open(winid, nil)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer w.CloseFiles()

	w.Addr("#1-+")
	b := make([]byte, 256)
	n, _ := w.Read("xdata", b)
	return string(b[:n-1]), nil
}
