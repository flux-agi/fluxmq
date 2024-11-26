package fluxmq

import (
	"log/slog"
)

type Options struct {
	clientName string
	host       string

	logger *slog.Logger
}

func WithHost(host string) ConnectionOpt {
	return func(c *Options) {
		c.host = host
	}
}

func WithLogger(logger *slog.Logger) ConnectionOpt {
	return func(c *Options) {
		c.logger = logger.WithGroup("fluxmq")
	}
}

func WithClientName(name string) ConnectionOpt {
	return func(c *Options) {
		c.clientName = name
	}
}
