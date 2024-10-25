package fluxnode

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/flux-agi/fluxmq/fluxmq"
)

type NodeStatus string

const (
	NodeStatusConnected NodeStatus = "CONNECTED"
	NodeStatusReady     NodeStatus = "READY"
	NodeStatusActive    NodeStatus = "ACTIVE"
	NodeStatusPaused    NodeStatus = "PAUSED"
	NodeStatusError     NodeStatus = "ERROR"
)

type Node[Settings any] struct {
	LifeCycle[Settings]

	logger     *slog.Logger
	ctx        context.Context
	alias      string
	connection *fluxmq.Connection

	nodeStatusLock sync.RWMutex
	nodeStatus     NodeStatus
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

	go func() {
		topic := fmt.Sprintf(topicStatusRequest, n.alias)
		ch, err := n.connection.Subscribe(n.ctx, topic)
		if err != nil {
			n.logger.Error("error for subscribe for req status")
			return
		}

		for {
			select {
			case <-n.ctx.Done():
				if err := n.connection.Unsubscribe(topicStatusRequest); err != nil {
					n.logger.Error("error for unsubscribe for req status")
				}
				return
			case <-ch:
				n.logger.Info("recived status request")
				for {
					n.nodeStatusLock.RLock()
					status := n.nodeStatus
					n.nodeStatusLock.RUnlock()

					if err := n.SetStatus(status); err != nil {
						n.logger.Error("error setting node status")
						time.Sleep(time.Second)
						continue
					}
					break
				}
			}
		}
	}()

	<-n.ctx.Done()

	return nil
}

// SetStatus for send status to manager
func (n *Node[T]) SetStatus(status NodeStatus) error {
	{
		n.nodeStatusLock.Lock()
		defer n.nodeStatusLock.Unlock()

		n.nodeStatus = status
	}

	topic := fmt.Sprintf(topicStatus, n.alias)
	if err := n.connection.Push(topic, []byte(status)); err != nil {
		return fmt.Errorf("push status: %w", err)
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
