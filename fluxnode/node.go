package fluxnode

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/flux-agi/fluxmq/fluxmq"
)

type TickSettings struct {
	IsInfinity bool          `json:"is_infinity"`
	Delay      time.Duration `json:"delay"`
}

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

	nodeLock       sync.RWMutex
	nodeStatus     NodeStatus
	tickSettingsCh chan TickSettings

	commonState []byte
}

// Create init flux service with settings
func Create[Settings any](ctx context.Context, logger *slog.Logger) *Node[Settings] {
	alias, ok := os.LookupEnv("NODE_ALIAS")
	if !ok {
		panic("env: 'NODE_ALIAS' must be set")
	}

	return &Node[Settings]{
		ctx:            ctx,
		alias:          alias,
		logger:         logger,
		tickSettingsCh: make(chan TickSettings),
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
		topicStatusReq := fmt.Sprintf(topicStatusRequest, n.alias)
		statusReqCh, err := n.connection.Subscribe(n.ctx, topicStatusReq)
		if err != nil {
			n.logger.Error("error for subscribe for req status")
			return
		}

		topicCommonStateReq := fmt.Sprintf(topicCommonState, n.alias)
		topicCommonStateCh, err := n.connection.Subscribe(n.ctx, topicCommonStateReq)
		if err != nil {
			n.logger.Error("error for subscribe for req state")
			return
		}

		for {
			select {
			case <-n.ctx.Done():
				if err := n.connection.Unsubscribe(topicStatusRequest); err != nil {
					n.logger.Error("error for unsubscribe for req status")
				}
				return
			case <-statusReqCh:
				n.logger.Info("recived status request")
				for {
					n.nodeLock.RLock()
					status := n.nodeStatus
					n.nodeLock.RUnlock()

					if err := n.SetStatus(status); err != nil {
						n.logger.Error("error setting node status")
						time.Sleep(time.Second)
						continue
					}
					break
				}

			case <-topicCommonStateCh:
				n.logger.Info("recived common state request")
				for {
					n.nodeLock.RLock()
					state := n.commonState
					n.nodeLock.RUnlock()

					t := fmt.Sprintf(topicPushCommonState, n.alias)
					if err := n.connection.Push(t, state); err != nil {
						n.logger.Error("error push node state")
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
		n.nodeLock.Lock()
		defer n.nodeLock.Unlock()

		n.nodeStatus = status
	}

	topic := fmt.Sprintf(topicStatus, n.alias)
	if err := n.connection.Push(topic, []byte(status)); err != nil {
		return fmt.Errorf("push status: %w", err)
	}
	return nil
}

// SetCommonState common state
func (n *Node[T]) SetCommonState(state any) error {
	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("marshal common state: %w", err)
	}

	n.nodeLock.Lock()
	defer n.nodeLock.Unlock()

	n.commonState = data
	return nil
}

// GetCommonState common state
func (n *Node[T]) GetCommonState() []byte {
	n.nodeLock.RLock()
	defer n.nodeLock.RUnlock()
	return n.commonState
}

// Alias of node
func (n *Node[S]) Alias() string {
	return n.alias
}

// Close connection of node
func (n *Node[S]) Close() error {
	close(n.tickSettingsCh)
	return n.connection.Close()
}
