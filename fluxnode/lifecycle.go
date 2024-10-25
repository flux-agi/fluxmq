package fluxnode

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/flux-agi/fluxmq/fluxmq"
)

type OnCommonFunc func() error
type OnConnectedFunc func(conn *fluxmq.Connection) error
type OnReadyFunc[T any] func(settings T) error
type OnTickFunc func(deltaTime time.Duration, timestamp time.Time)
type OnErrorFunc func(err error)

type LifeCycle[T any] struct {
	onConnected OnConnectedFunc
	onReady     OnReadyFunc[T]
	onTick      OnTickFunc
	onStart     OnCommonFunc
	onStop      OnCommonFunc
	onRestart   OnCommonFunc
	onError     OnErrorFunc
}

// OnConnected on connected event
func (n *Node[T]) OnConnected(f OnConnectedFunc) {
	n.onConnected = func(conn *fluxmq.Connection) error {
		if err := f(conn); err != nil {
			return err
		}

		if err := n.SetStatus(NodeStatusConnected); err != nil {
			return err
		}
		return nil
	}
}

// OnReady on ready event
func (n *Node[T]) OnReady(f OnReadyFunc[T]) {
	n.onReady = func(settings T) error {
		if err := f(settings); err != nil {
			return err
		}

		if err := n.SetStatus(NodeStatusReady); err != nil {
			return err
		}

		return nil
	}

	go func() {
		if n.connection == nil {
			for {
				time.Sleep(10 * time.Millisecond)
				if n.onConnected != nil {
					break
				}
			}
		}

		topic := fmt.Sprintf(topicOnReady, n.alias)
		ch, err := n.connection.Subscribe(n.ctx, topic)
		if err != nil {
			n.logger.Error("cannot subscribe to topic", slog.String("topic", topic), slog.Any("error", err))
			return
		}
		defer func() {
			if err := n.connection.Unsubscribe(topic); err != nil {
				n.logger.Error("cannot unsubscribe from topic", slog.String("topic", topic), slog.Any("error", err))
			}
		}()

		for {
			select {
			case <-n.ctx.Done():
				return
			case msg := <-ch:
				var data T
				if err := json.Unmarshal(msg.Payload, &data); err != nil {
					n.logger.Error("cannot unmarshal message", slog.Any("error", err))
					continue
				}

				if err := n.onReady(data); err != nil {
					slog.Info("onReady error", slog.Any("error", err))
					return
				}
			}
		}
	}()
}

// OnTick on tick event
func (n *Node[T]) OnTick(f OnTickFunc) {
	n.onTick = f

	// TODO: need implement on nats subscribtions
	go func() {
		for {
			select {
			case <-n.ctx.Done():
				return
			default:
				time.Sleep(250 * time.Millisecond)
				if n.onTick != nil {
					n.onTick(250*time.Millisecond, time.Now())
				}
			}
		}
	}()
}

// OnStart on start event
func (n *Node[T]) OnStart(f OnCommonFunc) {
	n.onStart = func() error {
		if err := f(); err != nil {
			return err
		}

		if err := n.SetStatus(NodeStatusActive); err != nil {
			return err
		}
		return nil
	}

	go func() {
		if n.connection == nil {
			for {
				time.Sleep(10 * time.Millisecond)
				if n.onConnected != nil {
					break
				}
			}
		}

		topic := fmt.Sprintf(topicOnStart, n.alias)
		ch, err := n.connection.Subscribe(n.ctx, topic)
		if err != nil {
			n.logger.Error("cannot subscribe to topic", slog.String("topic", topic), slog.Any("error", err))
			return
		}
		defer func() {
			if err := n.connection.Unsubscribe(topic); err != nil {
				n.logger.Error("cannot unsubscribe from topic", slog.String("topic", topic), slog.Any("error", err))
			}
		}()

		for {
			select {
			case <-n.ctx.Done():
				return
			case <-ch:
				if err := n.onStart(); err != nil {
					slog.Info("onStart error", slog.Any("error", err))
					return
				}
			}
		}
	}()
}

// OnStop on stop event
func (n *Node[T]) OnStop(f OnCommonFunc) {
	n.onStop = func() error {
		if err := f(); err != nil {
			return err
		}

		if err := n.SetStatus(NodeStatusPaused); err != nil {
			return err
		}
		return nil
	}

	go func() {
		if n.connection == nil {
			for {
				time.Sleep(10 * time.Millisecond)
				if n.onConnected != nil {
					break
				}
			}
		}

		topic := fmt.Sprintf(topicOnStop, n.alias)
		ch, err := n.connection.Subscribe(n.ctx, topic)
		if err != nil {
			n.logger.Error("cannot subscribe to topic", slog.String("topic", topic), slog.Any("error", err))
			return
		}
		defer func() {
			if err := n.connection.Unsubscribe(topic); err != nil {
				n.logger.Error("cannot unsubscribe from topic", slog.String("topic", topic), slog.Any("error", err))
			}
		}()

		for {
			select {
			case <-n.ctx.Done():
				return
			case <-ch:
				if err := n.onStop(); err != nil {
					slog.Info("onStop error", slog.Any("error", err))
					return
				}
			}
		}
	}()
}

// OnRestart on restart event
func (n *Node[T]) OnRestart(f OnCommonFunc) {

	n.onRestart = func() error {
		if err := n.SetStatus(NodeStatusPaused); err != nil {
			return err
		}

		if err := f(); err != nil {
			return err
		}

		if err := n.SetStatus(NodeStatusActive); err != nil {
			return err
		}
		return nil
	}

	go func() {
		if n.connection == nil {
			for {
				time.Sleep(10 * time.Millisecond)
				if n.onConnected != nil {
					break
				}
			}
		}

		topic := fmt.Sprintf(topicOnRestart, n.alias)
		ch, err := n.connection.Subscribe(n.ctx, topic)
		if err != nil {
			n.logger.Error("cannot subscribe to topic", slog.String("topic", topic), slog.Any("error", err))
			return
		}
		defer func() {
			if err := n.connection.Unsubscribe(topic); err != nil {
				n.logger.Error("cannot unsubscribe from topic", slog.String("topic", topic), slog.Any("error", err))
			}
		}()

		for {
			select {
			case <-n.ctx.Done():
				return
			case <-ch:
				if err := n.onRestart(); err != nil {
					slog.Info("onRestart error", slog.Any("error", err))
					return
				}
			}
		}
	}()
}

// OnError on error
func (n *Node[T]) OnError(f OnErrorFunc) {
	n.onError = f

	go func() {
		if n.connection == nil {
			for {
				time.Sleep(10 * time.Millisecond)
				if n.onConnected != nil {
					break
				}
			}
		}

		topic := fmt.Sprintf(topicOnError, "*")
		ch, err := n.connection.Subscribe(n.ctx, topic)
		if err != nil {
			n.logger.Error("cannot subscribe to topic", slog.String("topic", topic), slog.Any("error", err))
			return
		}
		defer func() {
			if err := n.connection.Unsubscribe(topic); err != nil {
				n.logger.Error("cannot unsubscribe from topic", slog.String("topic", topic), slog.Any("error", err))
			}
		}()

		for {
			select {
			case <-n.ctx.Done():
				return
			case msg := <-ch:
				e := errors.New(string(msg.Payload))
				n.onError(e)
			}
		}
	}()
}
