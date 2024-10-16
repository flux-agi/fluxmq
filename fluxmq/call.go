package fluxmq

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

// Call and wait response
func (c *Connection) Call(_ context.Context, subj string, data []byte) (*Msg, error) {
	resp, err := c.connection.Request(subj, data, time.Second)
	if err != nil {
		if errors.Is(err, nats.ErrNoResponders) {
			return nil, fmt.Errorf("no responders for subj: %s", subj)
		}
		return nil, err
	}

	msg := newDomainMsg(resp)

	return &msg, nil
}
