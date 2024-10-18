package fluxmq

import "github.com/nats-io/nats.go"

// Msg base mq message
type Msg struct {
	msg *nats.Msg

	Payload []byte `json:"payload"`
}

// newDomainMsg
func newDomainMsg(msg *nats.Msg) Msg {
	return Msg{
		msg:     msg,
		Payload: msg.Data,
	}
}

// NewMsg create new base mq message
func NewMsg(payload []byte) Msg {
	return Msg{Payload: payload}
}

// GetTopic for message
func (m *Msg) GetTopic() string { return m.msg.Subject }

// Respond response for request
func (m *Msg) Respond(data []byte) error {
	return m.msg.Respond(data)
}
