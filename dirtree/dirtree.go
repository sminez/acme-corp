// A directory tree viewer for acme.
//
// This make use of the Sam editing language in order to manipulate
// the acme `addr` file for the paired window, allowing us to set
// the position of the cursor, search back through the tree etc.
// The Same command language is _similar_ to [s]ed but has a bit more
// to it. See http://sam.cat-v.org/cheatsheet/ for a quick cheatsheet
// and look at http://doc.cat-v.org/bell_labs/structural_regexps/se.pdf
// for more on structural regular expressions.
//
// TODO: Only delete / insert the text that has changed as the result of
//       a user action. Probably would also want a `reset` action as well
//       in case the user messed up the state of the window.
package main

import (
	"flag"
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

// Formatters for laying out the directory tree
const (
	INDENT       = "  "
	DIRCOLLAPSED = "+ "
	DIREXPANDED  = "- "
	FILE         = "  "
	BUFFSIZE     = 1024
	SEPSIZE      = 2
)

// Represents a node in the filesystem (directory or file)
type node struct {
	name       string
	fullPath   string
	depth      int
	isDir      bool
	isExpanded bool
	showHidden bool
	contents   []*node
}

// The overall tree structure
type fileTree struct {
	w          *acme.Win
	root       string
	showHidden bool
	rootNodes  []*node
	nodeMap    map[string]*node
}

func main() {
	var root string

	flag.Usage = showUsage
	flag.Parse()
	args := flag.Args()

	switch len(args) {
	case 0:
		root, _ = os.Getwd()
	case 1:
		temp := path.Clean(args[0])
		if temp[0] != '/' {
			cwd, _ := os.Getwd()
			root = path.Join(cwd, temp)
		} else {
			root = temp
		}
	default:
		showUsage()
	}

	f := newFileTree(root)
	f.runEventLoop()
}

func showUsage() {
	fmt.Fprintf(os.Stderr, "usage: dirtree [path]\n")
	os.Exit(2)
}

// Set up a new dirtree window
func newFileTree(root string) *fileTree {
	f := fileTree{
		root:       root,
		showHidden: false,
		rootNodes:  []*node{},
		nodeMap:    make(map[string]*node),
	}

	f.rootNodes, _ = getNodes(f.root, 0)
	f.registerNodes(f.rootNodes)

	return &f
}

func escapePath(path string) string {
	return strings.Replace(path, " ", "\\ ", -1)
}

func getNodes(root string, depth int) ([]*node, error) {
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var nodes []*node

	for _, f := range fileInfo {
		name := f.Name()
		n := node{
			name:     name,
			fullPath: filepath.Join(root, name),
			depth:    depth,
			isDir:    f.IsDir(),
			contents: []*node{},
		}
		nodes = append(nodes, &n)
	}
	return nodes, nil
}

func (n *node) String() string {
	if !n.showHidden && n.name[0] == '.' {
		return ""
	}

	prefix := FILE
	if n.isDir {
		prefix = DIRCOLLAPSED

		if n.isExpanded {
			prefix = DIREXPANDED
		}
	}

	for i := 0; i < n.depth; i++ {
		prefix = prefix + INDENT
	}

	s := fmt.Sprintf("%s%s\n", prefix, n.name)
	if !(n.isDir && n.isExpanded) {
		return s
	}

	for _, m := range n.contents {
		m.showHidden = n.showHidden
		ms := m.String()
		if len(ms) > 0 {
			s = fmt.Sprintf("%s%s", s, ms)
		}
	}

	return s
}

// Send the filename to the plumber for opening. Depending on the filename
// this can still fail to open but at that point it is the responsibility
// of the user to write an appropriate plumbing rule.
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
		Data: []byte(escapePath(n.fullPath)),
	}

	return msg.Send(port)
}

// We clear/refetch the nodes on expand/collapse in order to allow
// the user to refresh the contents of a directory if, for example,
// they have created a new file.
func (f *fileTree) toggleDirectory(n *node) {
	if n.isExpanded {
		n.isExpanded = false
		for _, child := range n.contents {
			delete(f.nodeMap, child.fullPath)
		}
		n.contents = []*node{}
		return
	}

	var err error
	n.isExpanded = true
	n.contents, err = getNodes(n.fullPath, n.depth+1)
	if err != nil {
		f.w.Write("error", []byte(err.Error()))
		return
	}
	f.registerNodes(n.contents)
}

func (f *fileTree) String() string {
	// Wrap in parens for easy button3 to open the default acme file explorer
	s := fmt.Sprintf("(%s)\n\n", f.root)
	for _, n := range f.rootNodes {
		n.showHidden = f.showHidden
		ns := n.String()
		if len(ns) > 0 {
			s = fmt.Sprintf("%s%s", s, ns)
		}
	}
	return s
}

func (f *fileTree) registerNodes(ns []*node) {
	for _, n := range ns {
		f.nodeMap[n.fullPath] = n
	}
}

func (f *fileTree) resetRoot(root string) {
	f.root = root
	f.nodeMap = make(map[string]*node)
	f.rootNodes, _ = getNodes(f.root, 0)
	f.registerNodes(f.rootNodes)
	f.w.Name("+dirtree")
	f.redraw(nil)
}

// Infinite loop
func (f *fileTree) runEventLoop() {
	win, err := acme.New()
	if err != nil {
		fmt.Printf("Unable to initialise new acme window: %s\n", err)
		os.Exit(1)
	}

	win.Name("+dirtree")
	win.Write("tag", []byte("UpDir Hidden"))

	f.w = win
	f.redraw(nil)

	for e := range win.EventChan() {
		switch e.C2 {
		case 'x':
			// middle click in the tag
			switch strings.TrimSpace(string(e.Text)) {
			case "Del":
				win.Ctl("delete")

			case "Hidden":
				f.showHidden = !f.showHidden
				f.redraw(nil)

			case "UpDir":
				f.resetRoot(path.Dir(f.root))

			default:
				// Let acme handle it
				win.WriteEvent(e)
			}

		case 'L':
			// right click in body
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
				win.WriteEvent(e)
				continue
			}

			if n.isDir {
				f.toggleDirectory(n)
				f.redraw(e)
				continue
			}

			err = n.plumb()
			if err != nil {
				win.Write("error", []byte(err.Error()))
			}

		case 'X':
			// middle click in body
			path, err := f.getPath(e)
			if err != nil {
				f.w.Write("error", []byte(err.Error()))
				continue
			}

			if n := f.nodeMap[path]; n.isDir {
				f.resetRoot(n.fullPath)
			}

		default:
			win.WriteEvent(e)

		}
	}
}

func (f *fileTree) redraw(e *acme.Event) {
	f.w.Clear()
	f.w.Write("body", []byte(f.String()))

	if e != nil {
		// Keep the previous dot, moving to the start of the line
		f.w.Addr(fmt.Sprintf("#%d-1#0", e.OrigQ0))
	} else {
		f.w.Addr("1-1#0")
	}

	f.w.Ctl("dot=addr")
	f.w.Ctl("clean")
	f.w.Ctl("show")
}

func (f *fileTree) getPath(e *acme.Event) (string, error) {
	// Fetch the entire line from acme using addr
	// The acme address syntax here is going to the character at the
	// begining of the event selection text (#e.Orig0), jumping back
	// to the start of the line (-) and selecting to the end (+).
	f.w.Addr(fmt.Sprintf("#%d-+", e.OrigQ0))
	b := make([]byte, BUFFSIZE)
	n, _ := f.w.Read("xdata", b)
	line := string(b[:n-1])

	// Now that we have the line, get it's indentation level and walk
	// our way back up the window to get the rest of the path components.
	var ix int
	line = line[SEPSIZE:]

	for i := 0; i < len(line); i++ {
		if line[i] != ' ' {
			ix = i
			break
		}
	}

	indent := len(line[:ix]) / SEPSIZE
	clicked := line[ix:]

	p := []string{clicked}

	for i := indent - 1; i >= 0; i-- {
		spacer := ""
		for j := 0; j < i; j++ {
			spacer = spacer + INDENT
		}

		// Reverse search (-/regexp/) for the first line that is a directory
		// (starts with -/+) and is at the correct indentation level. Then
		// select the entire line.
		f.w.Addr(fmt.Sprintf(`-/[\-\+] %s[^ ]+/-+`, spacer))
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
