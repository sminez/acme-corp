package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"9fans.net/go/acme"
	"github.com/sminez/acme-corp/snoop-acme"
	"github.com/sminez/acme-corp/snoop-acme/afmt"
)

const (
	port            = 2009
	defaultSnoopTag = "Get fmton fmtoff dirtree"
)

// An AcmeSnooper snoops on acme events and listens for custom action requests over
// TCP. This allows for richer reuse of existing acme wrappers from acme.go
type AcmeSnooper struct {
	activeWindow int
	snoopWindow  *acme.Win
	listener     *snoop.Listener
	chLogEvents  chan acme.LogEvent
	formatOn     bool
}

// NewAcmeSnooper inits an acme snooper and grabs the /+snoop window so that we
// can send messages back to acme in a consistent way.
func NewAcmeSnooper() *AcmeSnooper {
	win, err := acme.New()
	if err != nil {
		log.Fatal(err)
	}
	pwd, _ := os.Getwd()
	win.Name(pwd + "/+snoop")
	win.Ctl("clean")
	win.Fprintf("tag", defaultSnoopTag)

	return &AcmeSnooper{
		activeWindow: -1,
		snoopWindow:  win,
		listener:     snoop.NewListener(port),
		chLogEvents:  make(chan acme.LogEvent),
	}
}

func (a *AcmeSnooper) writeToSnoopWindow(s string) {
	a.snoopWindow.Write("body", []byte(s))
}

// Snoop kicks off our local server and starts listening in on acme events.
func (a *AcmeSnooper) Snoop(chSignals chan os.Signal) {
	a.listener.RegisterHandler("active", a.activeHandler)
	a.listener.RegisterHandler("fmt", a.fmtHandler)

	go a.listener.HandleIncomingConnections()
	go a.tailLog()

	a.writeToSnoopWindow("acme snooper: Running")

	for {
		select {
		case e := <-a.chLogEvents:
			switch e.Op {
			case "":
				os.Exit(0) // acme was closed

			case "focus":
				a.activeWindow = e.ID

			case "put":
				if e.Name == "" || !a.formatOn {
					continue
				}

				for _, ft := range afmt.FTYPES {
					if ft.Matches(&e) {
						ft.Reformat(&e)
					}
				}

			default:
				// log.Printf("%s: %v\n", e.Op, e)
			}

		case <-chSignals:
			os.Exit(0)
		}
	}
}

func (a *AcmeSnooper) tailLog() {
	l, _ := acme.Log()
	for {
		e, _ := l.Read()
		a.chLogEvents <- e
	}
}

func (a *AcmeSnooper) fmtHandler(s string) (string, error) {
	switch s {
	case "on":
		a.formatOn = true
		a.writeToSnoopWindow("format on save: on")
		return "on", nil

	case "off":
		a.formatOn = false
		a.writeToSnoopWindow("format on save: off")
		return "off", nil

	default:
		return "", fmt.Errorf("'%s' is not a valid format directive", s)
	}
}

func (a *AcmeSnooper) activeHandler(s string) (string, error) {
	return string(a.activeWindow), nil
}

func main() {
	var progEndSignals = []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	}

	chSignals := make(chan os.Signal, 1)
	signal.Notify(chSignals, progEndSignals...)

	a := NewAcmeSnooper()
	a.Snoop(chSignals)
}
