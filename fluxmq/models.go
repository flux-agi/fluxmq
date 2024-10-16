package fluxmq

// Msg base mq message
type Msg struct {
	Payload []byte
}

// NewMsg create new base mq message
func NewMsg(payload []byte) Msg {
	return Msg{Payload: payload}
}
