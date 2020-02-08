// A directory tree viewer for acme.
//
// This make use of the Sam editing language in order to manipulate
// the acme `addr` file for the paired window, allowing us to set
// the position of the cursor, search back through the tree etc.
// The Sam command language is _similar_ to [s]ed but has a bit more
// to it. See http://sam.cat-v.org/cheatsheet/ for a quick cheatsheet
// and look at http://doc.cat-v.org/bell_labs/structural_regexps/se.pdf
// for more on structural regular expressions.
//
// TODO: - Only delete / insert the text that has changed as the result of
//       a user action. Probably would also want a `reset` action as well
//       in case the user messed up the state of the window.
//       - We don't follow the normal acme semantics for selecting text at
//       present: we map a button 3 click to the _line_ it was on not the
//       selected text. This is useful but potentially confusing...
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"9fans.net/go/acme"
	"9fans.net/go/plan9"
	"9fans.net/go/plumb"
)

const (
	INDENT       = "  "
	DIRCOLLAPSED = "+ "
	DIREXPANDED  = "- "
	FILE         = "  "
	BUFFSIZE     = 1024
	SEPSIZE      = 2
)

type node struct {
	name       string
	fullPath   string
	depth      int
	isDir      bool
	isExpanded bool
	isHidden   bool
	contents   []*node
}

type fileTree struct {
	w          *acme.Win
	root       string
	showHidden bool
	rootNodes  []*node
	nodeMap    map[string]*node
}

func main() {
	root, _ := os.Getwd()
	f := newFileTree(root)
	f.redraw(nil)
	f.runEventLoop()
}

func newFileTree(root string) *fileTree {
	win, err := acme.New()
	if err != nil {
		fmt.Printf("Unable to initialise new acme window: %s\n", err)
		os.Exit(1)
	}

	win.Name("+dirtree")
	win.Write("tag", []byte("UpDir Hidden"))
	rootNodes, _ := getNodes(root, 0)

	f := fileTree{
		w:          win,
		root:       root,
		showHidden: false,
		rootNodes:  rootNodes,
		nodeMap:    make(map[string]*node),
	}

	f.registerNodes(rootNodes)
	return &f
}

// Essentially just run 'ls' for the given root directory. We lazy list contents
// of directories as we expand them so we tag the nodes with their depth as they
// are created in order to track their path relative to the +dirtree window root.
func getNodes(root string, depth int) ([]*node, error) {
	var fileInfo []os.FileInfo
	var nodes []*node
	var err error

	if fileInfo, err = ioutil.ReadDir(root); err != nil {
		return nil, err
	}

	for _, f := range fileInfo {
		name := f.Name()
		n := node{
			name:     name,
			fullPath: filepath.Join(root, name),
			depth:    depth,
			isDir:    f.IsDir(),
			isHidden: name[0] == '.',
			contents: []*node{},
		}
		nodes = append(nodes, &n)
	}

	return nodes, nil
}

// Generate a string representation for this node and all child nodes if this is
// an expanded directory. Otherwise, just correctly indent this node for the tree.
func (n *node) stringifyRecursive(showHidden bool) string {
	if n.isHidden && !showHidden {
		return ""
	}

	prefix := FILE
	if n.isDir {
		if n.isExpanded {
			prefix = DIREXPANDED
		} else {
			prefix = DIRCOLLAPSED
		}
	}

	prefix = prefix + strings.Repeat(INDENT, n.depth)
	s := fmt.Sprintf("%s%s\n", prefix, n.name)

	if n.isDir && n.isExpanded {
		for _, m := range n.contents {
			s += m.stringifyRecursive(showHidden)
		}
	}

	return s
}

// Send the filename to the plumber for opening. Depending on the filename this
// can still fail to open but at that point it is the responsibility of the user
// to write an appropriate plumbing rule.
func (n *node) plumb() error {
	port, err := plumb.Open("send", plan9.OWRITE)
	if err != nil {
		return err
	}
	defer port.Close()

	msg := &plumb.Message{
		Src:  "dirtree",
		Dst:  "",
		Dir:  "/",
		Type: "text",
		Data: []byte(strings.Replace(n.fullPath, " ", "\\ ", -1)),
	}

	return msg.Send(port)
}

// We clear & refetch the nodes on expand/collapse in order to allow the user to
// refresh the contents of a directory.
func (f *fileTree) toggleDirectory(n *node) {
	switch n.isExpanded {
	case true:
		for _, child := range n.contents {
			delete(f.nodeMap, child.fullPath)
		}
		n.contents = []*node{}

	case false:
		var err error
		if n.contents, err = getNodes(n.fullPath, n.depth+1); err != nil {
			f.w.Write("error", []byte(err.Error()))
			return
		}
		f.registerNodes(n.contents)
	}

	n.isExpanded = !n.isExpanded
}

func (f *fileTree) registerNodes(ns []*node) {
	for _, n := range ns {
		f.nodeMap[n.fullPath] = n
	}
}

// Recursively stringify the current state of the entire tree. We also
// include the abspath to the current root node at the head of the
// window in order to make it easy to quickly perform other actions.
func (f *fileTree) String() string {
	var s string

	for _, n := range f.rootNodes {
		s += n.stringifyRecursive(f.showHidden)
	}

	return fmt.Sprintf("(%s)\n\n%s", f.root, s)
}

// Redraw the entire file tree window in its current state (currently this is
// incredibly inefficient). If this was triggered by an event (rather than an
// internal call from dirtree itself) we preserve the current 'dot' in acme,
// otherwise we set the 'dot' to the first line of the window.
func (f *fileTree) redraw(e *acme.Event) {
	f.w.Clear()
	f.w.Write("body", []byte(f.String()))

	if e != nil {
		f.w.Addr(fmt.Sprintf("#%d-1#0", e.OrigQ0))
	} else {
		f.w.Addr("1-1#0")
	}

	f.w.Ctl("dot=addr")
	f.w.Ctl("clean")
	f.w.Ctl("show")
}

func (f *fileTree) resetRoot(root string) {
	f.root = root
	f.nodeMap = make(map[string]*node)
	f.rootNodes, _ = getNodes(f.root, 0)
	f.registerNodes(f.rootNodes)
	f.w.Name("+dirtree")
	f.redraw(nil)
}

// loop over events we get from '+dirtree' until the user closes the window
func (f *fileTree) runEventLoop() {
	for e := range f.w.EventChan() {
		switch e.C2 {
		case 'x': // middle click in the tag
			switch strings.TrimSpace(string(e.Text)) {
			case "Del":
				f.w.Ctl("delete")

			case "Hidden":
				f.showHidden = !f.showHidden
				f.redraw(nil)

			case "UpDir":
				f.resetRoot(path.Dir(f.root))

			default:
				// Let acme handle it
				f.w.WriteEvent(e)
			}

		case 'X': // middle click in body
			path, err := f.getPath(e)
			if err != nil {
				f.w.Write("error", []byte(err.Error()))
				continue
			}

			if n := f.nodeMap[path]; n.isDir {
				f.resetRoot(n.fullPath)
			}

		case 'L': // right click in body
			path, err := f.getPath(e)
			if err != nil {
				f.w.Write("error", []byte(err.Error()))
				continue
			}

			n, ok := f.nodeMap[path]
			if !ok {
				// The path we generated didn't map to a known node so
				// this is most likely user entered text. Rather than
				// bail, see if acme knows what to do with it.
				f.w.WriteEvent(e)
				continue
			}

			if n.isDir {
				f.toggleDirectory(n)
				f.redraw(e)
				continue
			}

			if err = n.plumb(); err != nil {
				f.w.Write("error", []byte(err.Error()))
			}

		default:
			f.w.WriteEvent(e) // pass through

		}
	}
}

// Fetch the entire line from acme using addr. The acme address syntax here is
// going to the character at the begining of the event selection text (#e.Orig0),
// jumping back to the start of the line (-) and selecting to the end (+).
func (f *fileTree) getPath(e *acme.Event) (string, error) {
	f.w.Addr(fmt.Sprintf("#%d-+", e.OrigQ0))
	b := make([]byte, BUFFSIZE)
	n, _ := f.w.Read("xdata", b)
	line := string(b[:n-1])[SEPSIZE:]

	var ix int
	for i := 0; i < len(line); i++ {
		if line[i] != ' ' {
			ix = i
			break
		}
	}

	indent := len(line[:ix]) / SEPSIZE
	p := []string{line[ix:]}

	// Now that we have the line, get it's indentation level and walk
	// our way back up the window to get the rest of the path components.
	for i := indent - 1; i >= 0; i-- {
		// Reverse search (-/regexp/) for the first line that is a directory
		// (starts with -/+) and is at the correct indentation level. Then
		// select the entire line.
		f.w.Addr(fmt.Sprintf(`-/[\-\+] %s[^ ]+/-+`, strings.Repeat(INDENT, i)))
		b := make([]byte, BUFFSIZE)
		n, _ := f.w.Read("xdata", b)
		comp := strings.TrimSpace(string(b[:n-1])[SEPSIZE:])
		p = append(p, comp)
	}

	// Reverse to get everything in path order
	comps := []string{f.root}
	for i := len(p) - 1; i >= 0; i-- {
		comps = append(comps, p[i])
	}

	return path.Join(comps...), nil
}
