package fluxmq

// Respond [preferred for use]  answer for request
func (c *Connection) Respond(msg *Msg, data []byte) error {
	return msg.Respond(data)
}
