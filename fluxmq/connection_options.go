package fluxmq

import (
	"log/slog"
)

// WithHost set connection host
func WithHost(host string) ConnectionOpt {
	return func(c *Connection) {
		c.host = host
	}
}

func WithLogger(logger *slog.Logger) ConnectionOpt {
	return func(c *Connection) {
		c.logger = logger.WithGroup("fluxmq")
	}
}
