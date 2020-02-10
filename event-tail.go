package main

import (
	"fmt"
	"log"

	"9fans.net/go/acme"
)

const winID = 2 // just testing for now using the snooper window since it's always there

// -- notes (go doc 9fans.net/go/acme Event) --
// * acme.Event C1 : 'K' for keyboard input, 'M' for mouse clicks
// * Intercepting keyboard events does _not_ seem to prevent them going to the window (unlike
//   mouse events)
// * For simple mouse actions, it looks like event.Text is enough to get what was clicked on.
//   For more complicated chording / large selections you need to look at the other event fields
//   and the go fetch them from the window yourself using the addr.
// * The win.EventLoop function isn't going to be right for this, it's only allowing you to
//   handle mouse events and then on top of that it seems to use reflection to dynamically pull
//   off methods from the EventHandler type you create in order to do anything.
//   * It is also accessing a lot of private stuff on acme.Win so I might need to reimplement some
//     stuff to get this to work...

func main() {
	w, err := acme.Open(winID, nil)
	if err != nil {
		log.Print(err)
	}
	defer w.CloseFiles()

	for e := range w.EventChan() {
		fmt.Printf("%#v\n", e)
		//w.WriteEvent(e)  // Commenting this out doesn't seem to prevent keys being sent
	}
}
