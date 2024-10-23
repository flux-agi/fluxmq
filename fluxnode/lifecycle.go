package fluxnode

import (
	"time"

	"github.com/flux-agi/fluxmq/fluxmq"
)

type OnCommonFunc func() error
type OnConnectedFunc func(conn *fluxmq.Connection) error
type OnReadyFunc[T any] func(settings T) error
type OnTickFunc func(deltaTime time.Duration, timestamp time.Time)
type OnErrorFunc func(err error)

type LifeCycle[T any] struct {
	onConnected   OnConnectedFunc
	onInitialized OnCommonFunc
	onReady       OnReadyFunc[T]
	onTick        OnTickFunc
	onStart       OnCommonFunc
	onStop        OnCommonFunc
	onRestart     OnCommonFunc
	onError       OnErrorFunc
}

func (n *Node[T]) OnConnected(f OnConnectedFunc) { n.onConnected = f }
func (n *Node[T]) OnInitialized(f OnCommonFunc)  { n.onInitialized = f }
func (n *Node[T]) OnReady(f OnReadyFunc[T])      { n.onReady = f }
func (n *Node[T]) OnTick(f OnTickFunc)           { n.onTick = f }
func (n *Node[T]) OnStart(f OnCommonFunc)        { n.onStart = f }
func (n *Node[T]) OnStop(f OnCommonFunc)         { n.onStop = f }
func (n *Node[T]) OnRestart(f OnCommonFunc)      { n.onRestart = f }
func (n *Node[T]) OnError(f OnErrorFunc)         { n.onError = f }
