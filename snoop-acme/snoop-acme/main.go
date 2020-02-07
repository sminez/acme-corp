package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sminez/acme-corp/snoop-acme"
)

func main() {
	var progEndSignals = []os.Signal{
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
	}

	chSignals := make(chan os.Signal, 1)
	signal.Notify(chSignals, progEndSignals...)

	a := snoop.NewAcmeSnooper()
	a.Snoop(chSignals)
}
