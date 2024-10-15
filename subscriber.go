package fluxmq

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	zeromq "github.com/pebbe/zmq4"
)

// Subscriber is subscriber
type Subscriber struct {
	socket *zeromq.Socket
}

// CreateSub create new subscriber
func CreateSub(topic string) (*Subscriber, error) {
	socket, err := zeromq.NewSocket(zeromq.SUB)
	if err != nil {
		return nil, fmt.Errorf("cannot create socket: %s", err)
	}

	if err = socket.Connect(baseSubAddr); err != nil {
		return nil, fmt.Errorf("cannot bind socket to %s: %s", basePubAddr, err)
	}

	if err := socket.SetSubscribe(topic); err != nil {
		return nil, fmt.Errorf("cannot subscribe to %s: %s", topic, err)
	}

	return &Subscriber{
		socket: socket,
	}, nil
}

// Subscribe receive message
func (p *Subscriber) Subscribe(ctx context.Context) (chan []byte, error) {
	ch := make(chan []byte)
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				data, err := p.socket.Recv(0)
				if err != nil {
					slog.Error("cannot receive message: %s", slog.Any("err", err))
					return
				}
				fmt.Printf("receive message: %s\n", data)
				ch <- []byte(data)
			}

		}
	}(ctx)
	return ch, nil
}

// Close subscriber
func (p *Subscriber) Close() error {
	var err error

	if err := p.socket.Disconnect(baseSubAddr); err != nil {
		err = errors.Join(fmt.Errorf("cannot disconnect from %s: %s", baseSubAddr, err), err)
	}

	if err := p.socket.Close(); err != nil {
		err = errors.Join(fmt.Errorf("cannot close socket: %s", err), err)
	}

	return err
}
