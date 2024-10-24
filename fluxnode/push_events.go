package fluxnode

import "fmt"

// requestConfig request config for service
func (n *Node[T]) requestConfig() error {
	topic := fmt.Sprintf(topicRequestConfig, n.alias)
	if err := n.connection.Push(topic, []byte(n.alias)); err != nil {
		return err
	}
	return nil
}

// PushError push error
func (n *Node[T]) PushError(err error) error {
	topic := fmt.Sprintf(topicOnError, n.alias)
	if err := n.connection.Push(topic, []byte(err.Error())); err != nil {
		return err
	}
	return nil
}
