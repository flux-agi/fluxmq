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
			if err := c.Unsubscribe(s.Subject); err != nil {
				c.logger.Error("error unsubscribe",
					slog.String("topic", topic),
					slog.Any("error", err),
				)
			}

			close(ch)

			c.logger.Info("closing subscriber channel",
				slog.String("topic", topic),
			)
		}()

		<-ctx.Done()
	}()

	if _, ok := c.subscriptions.Load(topic); !ok {
		c.subscriptions.Store(topic, s)
	} else {
		c.logger.Warn("subscription already exists", slog.String("topic", topic))
	}

	return ch, nil
}

// Unsubscribe from topic
func (c *Connection) Unsubscribe(topic string) error {
	sub, ok := c.subscriptions.Load(topic)
	if !ok {
		return fmt.Errorf("subscription for topic %s not found", topic)
	}

	defer c.subscriptions.Delete(topic)

	if err := sub.(*nats.Subscription).Unsubscribe(); err != nil {
		return fmt.Errorf("failed to unsubscribe: %w", err)
	}
	return nil
}
