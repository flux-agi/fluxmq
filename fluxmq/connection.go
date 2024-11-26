package fluxmq

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/nats-io/nats.go"
)

// ConnectionOpt is connection option
type ConnectionOpt func(*Options)
type ClientSubscription map[string]*nats.Subscription

// Connection is connection
type Connection struct {
	logger        *slog.Logger
	connection    *nats.Conn
	subscriptions *sync.Map
}

// Connect create new connection to fluxMQ
func Connect(opts ...ConnectionOpt) (*Connection, error) {
	o := &Options{
		clientName: "",
		host:       nats.DefaultURL,
		logger:     slog.Default().WithGroup("fluxmq"),
	}

	for _, opt := range opts {
		opt(o)
	}

	connection, err := nats.Connect(
		o.host,
		nats.Name(o.clientName),
	)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to %s: %s", o.host, err)
	}

	return &Connection{
		logger:        o.logger,
		connection:    connection,
		subscriptions: new(sync.Map),
	}, nil
}

// Close connection
func (c *Connection) Close() error {
	c.connection.Close()
	return nil
}
