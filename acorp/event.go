package acorp

import (
	"fmt"

	"9fans.net/go/acme"
)

// A handler fucntion that processes an Acme event and takes an action. Passthrough must be explicitly
// carried out by the handler function itself.
type handler = func(*acme.Win, *acme.Event, func() error) error

// An EventFilter takes hold of an acme window's event file and passes all events
// it sees through a set of filter functions if they are defined. Unmatched events
// are passed through to acme.
type EventFilter struct {
	complete           bool
	KeyboardInputBody  handler
	KeyboardDeleteBody handler
	KeyboardInputTag   handler
	KeyboardDeleteTag  handler
	Mouse2Body         handler
	Mouse3Body         handler
	Mouse2Tag          handler
	Mouse3Tag          handler
}

func (ef *EventFilter) markComplete() error {
	ef.complete = true
	return nil
}

func (ef *EventFilter) applyOrPassthrough(f handler, w *acme.Win, e *acme.Event) error {
	if f != nil {
		return f(w, e, ef.markComplete)
	}

	return w.WriteEvent(e)
}

// Filter runs the event filter, releasing the window event file on the first error encountered
func (ef *EventFilter) Filter(w *acme.Win) error {
	for e := range w.EventChan() {
		if err := ef.filterSingle(w, e); err != nil {
			return err
		}

		if ef.complete {
			return nil
		}
	}

	return fmt.Errorf("lost event channel")
}

// Currently dropping E and F events that are generated by writes from other programs to the acme
// control files.
func (ef *EventFilter) filterSingle(w *acme.Win, e *acme.Event) error {
	switch e.C1 {
	case 'K':
		switch e.C2 {
		case 'I':
			ef.applyOrPassthrough(ef.KeyboardInputBody, w, e)
		case 'D':
			ef.applyOrPassthrough(ef.KeyboardDeleteBody, w, e)
		case 'i':
			ef.applyOrPassthrough(ef.KeyboardInputTag, w, e)
		case 'd':
			ef.applyOrPassthrough(ef.KeyboardDeleteTag, w, e)
		}
		return nil

	case 'M':
		switch e.C2 {
		case 'X':
			return ef.applyOrPassthrough(ef.Mouse2Body, w, e)
		case 'L':
			return ef.applyOrPassthrough(ef.Mouse3Body, w, e)
		case 'x':
			return ef.applyOrPassthrough(ef.Mouse2Tag, w, e)
		case 'l':
			return ef.applyOrPassthrough(ef.Mouse3Tag, w, e)
		}
	}

	w.WriteEvent(e)
	return nil
}
