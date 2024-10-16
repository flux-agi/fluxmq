package transport

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

// Call and wait response
func (c *Connection) Call(_ context.Context, topic string, data []byte) (*Msg, error) {
	resp, err := c.connection.Request(topic, data, time.Second)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return nil, fmt.Errorf("no responders for topic: %s", topic)
		}
		return nil, err
	}

	msg := newDomainMsg(resp)

	return &msg, nil
}
