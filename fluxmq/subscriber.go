package fluxmq

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go"
)

// Subscribe receive message.
//
// Сhan will close if ctx is Done, cancel it by "context.WithCancel()"
func (c *Connection) Subscribe(ctx context.Context, topic string) (chan Msg, error) {
	var ch = make(chan Msg)
	s, err := c.connection.Subscribe(topic, func(msg *nats.Msg) {
		c.logger.Info("recv msg", slog.String("data", string(msg.Data)))
		ch <- newDomainMsg(msg)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to topic %s: %w", topic, err)
	}

	s.SetClosedHandler(func(subject string) {
		c.logger.Info("subscriber closed", slog.String("subject", subject))
	})

	go func() {
		defer func() {
			close(ch)
			c.logger.Info("closing subscriber channel",
				slog.String("topic", topic),
			)
		}()

		select {
		case <-ctx.Done():
			if err := s.Unsubscribe(); err != nil {
				c.logger.Error("error unsubscribe",
					slog.String("topic", topic),
					slog.Any("error", err),
				)
			}
		}
	}()
	return ch, nil
}
