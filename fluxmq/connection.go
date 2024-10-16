package fluxmq

import (
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"
)

// ConnectionOpt is connection option
type ConnectionOpt func(*Connection)

// Connection is connection
type Connection struct {
	connection *nats.Conn

	host   string
	logger *slog.Logger
}

// Connect create new connection to fluxMQ
func Connect(opts ...ConnectionOpt) (*Connection, error) {
	var c Connection
	for _, opt := range opts {
		opt(&c)
	}

	setDefaultsIfNeeded(&c)

	connection, err := nats.Connect(c.host)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to %s: %s", c.host, err)
	}

	c.connection = connection

	return &c, nil
}

// setDefaultsIfNeeded set default values if needed for make connection
func setDefaultsIfNeeded(c *Connection) {
	if c.host == "" {
		c.host = nats.DefaultURL
	}

	if c.logger == nil {
		logger := slog.Default()
		c.logger = logger.WithGroup("fluxmq")
	}
}

// Close connection
func (c *Connection) Close() error {
	c.connection.Close()
	return nil
}
