package acorp

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"

	"9fans.net/go/acme"
)

const (
	activeWindowQuery = "active / .\n"
	snooperAddr       = "127.0.0.1:2009"
)

func winIDFromSnooper() (string, error) {
	conn, _ := net.Dial("tcp", snooperAddr)
	fmt.Fprintf(conn, activeWindowQuery)
	message, _ := bufio.NewReader(conn).ReadString('\n')
	if message == "-1" {
		return "", fmt.Errorf("unable to determine current window ID")
	}
	return message, nil
}

// GetCurrentWindow finds the current active window in acme, using the snooper if this is not called from
// inside of an acme window directly.
func GetCurrentWindow() (*acme.Win, error) {
	var err error

	winStr := os.Getenv("winid")
	if len(winStr) == 0 {
		winStr, err = winIDFromSnooper()
		if err != nil {
			return nil, fmt.Errorf("unable to determine current acme window")
		}
	}

	winID, err := strconv.Atoi(winStr)
	if err != nil {
		return nil, fmt.Errorf("non numeric winid: %s", winStr)
	}

	w, err := acme.Open(winID, nil)
	if err != nil {
		return nil, err
	}

	return w, nil
}
