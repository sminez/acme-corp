package acorp

import (
	"fmt"

	"9fans.net/go/acme"
)

// A KeyFilter can chose to intercept a given KeyEvent or pass it through to the
// underlying acme window. If intercepted, a pointer to the current window is
// included on the KeyEvent to allow for manipulating the window in response to
// the event.
type KeyFilter func(k KeyEvent) KeyEvent

// A KeyEvent is a wrapper around the raw runes seen from the acme event stream
// in order to simplify writing interceptors.
type KeyEvent struct {
	// The raw run we saw from acme
	r rune
	// Was the control key held as part of this event?
	Ctrl bool
	// The underlying acme window
	Win *acme.Win
}

func (k *KeyEvent) String() string {
	if k.r <= 26 {
		return fmt.Sprintf("C-%c", k.r+96)
	}
	return string(k.r)
}

// InterceptKeys tails the `acme/$winid/event` file for a given window and
// filters all Keyboard inputs through `fn`, which can decide to handle the
// envent itself of return the event so that it can be passed through to the
// underlying window.The returned channel can be used to signal shutdown.
func InterceptKeys(w *acme.Win, fn KeyFilter) chan struct{} {
	stop := make(chan struct{})
	return stop
}
