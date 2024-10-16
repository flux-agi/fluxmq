package fluxmq

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"
)

// Subscriber is subscriber
type Subscriber struct {
	conn   *nats.Conn
	logger *slog.Logger

	topic string
}

// CreateSub create new subscriber
func (c *Connection) CreateSub(topic string) (*Subscriber, error) {
	return &Subscriber{
		conn:   c.connection,
		logger: c.logger,
		topic:  topic,
	}, nil
}

// Recv receive message.
//
// Сhan will close if ctx is Done, cancel it by "context.WithCancel()"
func (p *Subscriber) Recv(ctx context.Context) (chan Msg, error) {
	var ch = make(chan Msg)
	s, err := p.conn.Subscribe(p.topic, func(msg *nats.Msg) {
		p.logger.Info("recv msg", slog.String("data", string(msg.Data)))
		ch <- NewMsg(msg.Data)
	})

	go func() {
		defer func() {
			close(ch)
			p.logger.Info("closing subscriber channel",
				slog.String("topic", p.topic),
			)
		}()

		select {
		case <-ctx.Done():
			if err := s.Unsubscribe(); err != nil {
				p.logger.Error("error unsubscribe",
					slog.String("topic", p.topic),
					slog.Any("error", err),
				)
			}
		}
	}()

	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to %s: %w", p.topic, err)
	}
	return ch, nil
}
