package pkg

import (
	"log/slog"

	"github.com/nats-io/nats.go"
)

// Publisher is publisher
type Publisher struct {
	conn   *nats.Conn
	logger *slog.Logger
}

// CreatePub create new publisher
func (c *Connection) CreatePub() (*Publisher, error) {
	return &Publisher{
		conn:   c.connection,
		logger: c.logger,
	}, nil
}

// Push message to topic
func (p *Publisher) Push(topic string, message []byte) error {
	defer p.logger.Info("send message",
		slog.String("topic", topic),
		slog.String("message", string(message)),
	)
	return p.conn.Publish(topic, message)
}
