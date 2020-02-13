package snoop

import (
	"fmt"
	"log"
	"os"

	"9fans.net/go/acme"
)

// An AcmeSnooper snoops on acme events and listens for custom action requests over
// TCP. This allows for richer reuse of existing acme wrappers from acme.go
type AcmeSnooper struct {
	win         *acme.Win
	listener    *Listener
	chLogEvents chan acme.LogEvent
	active      int
	formatOn    bool
	debug       bool
}

// NewAcmeSnooper inits an acme snooper and grabs the /+snoop window so that we
// can send messages back to acme in a consistent way.
func NewAcmeSnooper(debug bool) *AcmeSnooper {
	win, err := acme.New()
	if err != nil {
		log.Fatal(err)
	}
	win.Name("+snoop")
	win.Ctl("clean")
	win.Write("tag", []byte(defaultSnoopTag))

	return &AcmeSnooper{
		win:         win,
		listener:    NewListener(tcpPort),
		chLogEvents: make(chan acme.LogEvent),
		active:      -1,
		formatOn:    false,
		debug:       debug,
	}
}

func (a *AcmeSnooper) logf(s string, args ...interface{}) {
	a.win.Write("body", []byte(prompt+fmt.Sprintf(s, args...)))
	a.win.Ctl("clean")
}

func (a *AcmeSnooper) errorf(s string, args ...interface{}) {
	a.win.Write("errors", []byte(fmt.Sprintf(s, args...)))
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
		a.logf("format on save: enabled\n")
		return "on", nil

	case "off":
		a.formatOn = false
		a.logf("format on save: disabled\n")
		return "off", nil

	default:
		return "", fmt.Errorf("'%s' is not a valid format directive", s)
	}
}

func (a *AcmeSnooper) activeHandler(s string) (string, error) {
	return fmt.Sprintf("%d", a.active), nil
}

// Snoop kicks off our local server and starts listening in on acme events.
func (a *AcmeSnooper) Snoop(chSignals chan os.Signal) {
	a.listener.Register("active", a.activeHandler)
	a.listener.Register("fmt", a.fmtHandler)

	go a.listener.HandleIncomingConnections()
	go a.tailLog()

	a.win.Write("body", []byte("-- acme corp --\n"))
	a.logf("snooper now running...\n")

	for {
		select {
		case e := <-a.chLogEvents:
			switch e.Op {
			case "":
				os.Exit(0) // acme was closed

			case "focus":
				a.active = e.ID

			case "put":
				if a.formatOn && len(e.Name) > 0 {
					for _, ft := range formatableTypes {
						if ft.Matches(&e) {
							s := ft.Reformat(&e)
							if len(s) > 0 {
								a.errorf(s)
							}
							break
						}
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
