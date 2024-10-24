package fluxnode

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/flux-agi/fluxmq/fluxmq"
)

type Node[Settings any] struct {
	LifeCycle[Settings]

	logger     *slog.Logger
	ctx        context.Context
	alias      string
	connection *fluxmq.Connection
}

// Create init flux service with settings
func Create[Settings any](ctx context.Context, alias string, logger *slog.Logger) *Node[Settings] {
	return &Node[Settings]{
		ctx:    ctx,
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

	if conn == nil {
		return fmt.Errorf("connection is nil")
	}

	n.connection = conn

	if n.onConnected != nil {
		if err := n.onConnected(conn); err != nil {
			return fmt.Errorf("onConnected: %w", err)
		}
	}

	if err := n.requestConfig(); err != nil {
		return fmt.Errorf("requestConfig: %w", err)
	}

	<-n.ctx.Done()

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
