package fluxnode

import (
	"fmt"
	"log/slog"

	"github.com/flux-agi/fluxmq/fluxmq"
)

type Node[Settings any] struct {
	LifeCycle[Settings]

	alias      string
	logger     *slog.Logger
	connection *fluxmq.Connection
}

// Create init flux service with settings
func Create[Settings any](alias string, logger *slog.Logger) *Node[Settings] {
	return &Node[Settings]{
		alias:  alias,
		logger: logger,
	}
}

// Run node watcher
func (n *Node[T]) Run(connOpts ...fluxmq.ConnectionOpt) error {
	conn, err := fluxmq.Connect(connOpts...)
	if err != nil {
		return fmt.Errorf("connect to fluxmq server: %w", err)
	}

	n.connection = conn

	if err := n.onConnected(conn); err != nil {
		n.logger.Error("onConnected err", slog.Any("err", err))
	}

	return nil
}

// Alias of node
func (n *Node[S]) Alias() string {
	return n.alias
}

// Close connection of node
func (n *Node[S]) Close() error {
	return n.connection.Close()
}
