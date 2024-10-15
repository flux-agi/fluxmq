package fluxmq

import (
	"fmt"
	"strings"

	zeromq "github.com/pebbe/zmq4"
)

// Publisher is publisher
type Publisher struct {
	socket *zeromq.Socket
}

// CreatePub create new publisher
func CreatePub() (*Publisher, error) {
	socket, err := zeromq.NewSocket(zeromq.PUB)
	if err != nil {
		return nil, fmt.Errorf("cannot create socket: %s", err)
	}

	if err = socket.Bind(basePubAddr); err != nil {
		return nil, fmt.Errorf("cannot bind socket to %s: %s", basePubAddr, err)
	}

	return &Publisher{
		socket: socket,
	}, nil
}

// Send message to topic
func (p *Publisher) Send(topic, message string) error {
	if strings.Contains(topic, "/") ||
		strings.Contains(topic, "_") {
		return fmt.Errorf("invalid topic: %s; topic should not contain '/' or '_'", topic)
	}
	_, err := p.socket.Send(fmt.Sprintf("%s %s", topic, message), 0)
	return err
}

// Close publisher
func (p *Publisher) Close() error {
	return p.socket.Close()
}
