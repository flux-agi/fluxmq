package fluxmq

import (
	"log/slog"
)

// Push message to topic
func (c *Connection) Push(topic string, message []byte) error {
	defer c.logger.Info("send message",
		slog.String("topic", topic),
		slog.String("message", string(message)),
	)
	return c.connection.Publish(topic, message)
}
