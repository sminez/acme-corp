package snoop

// This file implements a simple TCP based protocol allowing clients to submit
// messages of the form '<route> / <content>' to the snooper for processing.
// Messages at this level are all arbitrary strings with parsing and processing
// being done by the MessageHandlers themselves. Invalid messages will have an
// error message returned to them but the format of that error is not
// guaranteed by the server.

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// A MessageHandler is a function that knows how to parse a given message type
type MessageHandler func(s string) (string, error)

// A Message is a simple RPC message format to allow for very simple scripts
// that pass simple string messages to the snooper for it to process. Ideally
// this should be exposed as a 9fs file system in the same way as acme itself
// but for now this will have to do.
type Message struct {
	route   string
	content string
}

// NewMessage constructs a new Message from a string that we received over the
// network. At this stage we are not guaranteed that this message is valid,
// only that we were able to split it into a handler and content.
func NewMessage(s string) (*Message, error) {
	sections := strings.SplitN(s, "/", 2)
	if len(sections) != 2 {
		return nil, fmt.Errorf("Invalid message '%s'", s)
	}
	return &Message{
		route:   strings.TrimSpace(sections[0]),
		content: strings.TrimSpace(sections[1]),
	}, nil
}

// A Listener runs the event loop that owns our TCP socket and routes incoming
// messages to their relevant handlers.
type Listener struct {
	handlers map[string]MessageHandler
	port     int
}

// NewListener initialises a new Listener without any handlers
func NewListener(port int) *Listener {
	return &Listener{
		handlers: make(map[string]MessageHandler),
		port:     port,
	}
}

// HandleIncomingConnections binds to a tcp socket and serves handler responses
// for incoming connections. Runs in a goroutine.
func (l *Listener) HandleIncomingConnections() {
	s, _ := net.Listen("tcp", fmt.Sprintf("localhost:%d", l.port))
	for {
		// silently dropping failed incoming connections
		conn, _ := s.Accept()
		go l.handleConnection(conn)
	}
}

// Register registers a new message handler with a given route
func (l *Listener) Register(route string, handler MessageHandler) {
	l.handlers[route] = handler
}

// Runs in a goroutine per incoming connection
func (l *Listener) handleConnection(conn net.Conn) {
	s, _ := bufio.NewReader(conn).ReadString('\n')
	defer conn.Close()

	msg, err := NewMessage(s)
	if err != nil {
		conn.Write([]byte(err.Error()))
		return
	}

	handler, ok := l.handlers[msg.route]
	if !ok {
		conn.Write([]byte(fmt.Sprintf("'%s' is not a known handler", msg.route)))
		return
	}

	resp, err := handler(msg.content)
	if err != nil {
		conn.Write([]byte(err.Error()))
		return
	}
	conn.Write([]byte(resp))
}
