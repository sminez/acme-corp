/*
`punt` allows you to edit the contents of an acme window in an external editor.

It creates a new temp file to work with rather than opening the underlying source
file again inside the new editor and changes are written to the acme window itself
not the underlying file.

File extensions are parsed when creating thetemp file so that things like plugins
and syntax highlighting trigger correctly in the spawned editor. At present, the
launch of the editor is done via spawning a tilix terminal session first so that
terminal based editors (vim, emacs, nano...) can be used. If the editor in
question is GUI based then you will see an additional terminal window as well.
*/
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"9fans.net/go/acme"
	"github.com/fsnotify/fsnotify"
)

var (
	terminal      = "tilix"
	defaultEditor = "nvim"
	tmpName       = "acme-punt"
)

// Default to `defaultEditor` running in `terminal` or use the user provided details
func parseArgs() (string, bool) {
	isGUI := flag.Bool("g", false, "The called program is a GUI")
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		return defaultEditor, false
	}

	return args[0], *isGUI
}

func main() {
	editor, isGUI := parseArgs()

	winStr := os.Getenv("winid")
	if len(winStr) == 0 {
		log.Fatal("Need to be called from inside acme")
	}

	winID, err := strconv.Atoi(winStr)
	if err != nil {
		log.Fatalf("Non numeric winid: %s\n", winStr)
	}

	w, err := acme.Open(winID, nil)
	if err != nil {
		log.Print(err)
	}
	defer w.CloseFiles()

	winInfo, err := acme.Windows()
	if err != nil {
		log.Print(err)
		return
	}

	for _, i := range winInfo {
		if i.ID == winID {
			comps := strings.SplitAfter(i.Name, ".")
			if n := len(comps); n > 1 {
				suffix := comps[n-1]
				tmpName = fmt.Sprintf("acme-punt.*.%s", suffix)
			}
			break
		}
	}

	w.Addr(",")
	body, err := w.ReadAll("data")
	if err != nil {
		log.Print(err)
		return
	}

	// Set up a temp file, copy the contents of the current window to it and
	// then open it up in that editor.
	tmpFile, err := ioutil.TempFile("", tmpName)
	if err != nil {
		log.Print(err)
		return
	}

	if _, err := tmpFile.Write(body); err != nil {
		tmpFile.Close()
		log.Fatal(err)
	}
	tmpFile.Close()

	// Kick off a file watcher to track when the file is saved
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Print(err)
		return
	}

	go func() {
		defer os.Remove(tmpFile.Name()) // clean up
		defer watcher.Close()

		if err := watcher.Add(tmpFile.Name()); err != nil {
			log.Printf("fsnotify-error: %s\n", err)
			return
		}

		for {
			select {
			// watch for events
			case event := <-watcher.Events:
				// log.Printf("fsnotify: %#v\n", event)

				if event.Op&fsnotify.Write == fsnotify.Write {
					// Read the tempfile contents back in and replace the
					// window content with them.
					// TODO: Stash cursor position? Can't guarentee that we'll
					// still have a similar structure when we come back...
					edited, err := ioutil.ReadFile(tmpFile.Name())
					if err != nil {
						log.Print(err)
						return
					}

					w.Clear()
					w.Write("data", edited)
					os.Remove(tmpFile.Name())
					return
				}

				if event.Op&fsnotify.Remove == fsnotify.Remove {
					// If the file was removed then kill the goro
					return
				}

			case err := <-watcher.Errors:
				log.Printf("fsnotify-error: %s\n", err)
			}
		}
	}()

	// spawn the editor
	if isGUI {
		_, err = exec.Command(editor, tmpFile.Name()).CombinedOutput()
	} else {
		_, err = exec.Command(terminal, "-e", editor, tmpFile.Name()).CombinedOutput()
	}
	if err != nil {
		log.Print(err)
	}
}
