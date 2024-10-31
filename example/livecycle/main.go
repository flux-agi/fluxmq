package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"github.com/flux-agi/fluxmq/fluxmq"
	"github.com/flux-agi/fluxmq/fluxnode"
)

const serviceAlias = "fluxnode"

type settings struct {
}

func main() {
	ctx := context.Background()
	logger := slog.Default()

	os.Setenv("NODE_ALIAS", serviceAlias)

	node := fluxnode.Create[settings](ctx, logger)

	node.OnConnected(func(conn *fluxmq.Connection) error {
		logger.Info("connected")

		go func() {
			for {
				time.Sleep(time.Second)

				tickdata, _ := json.Marshal(fluxnode.TickSettings{
					IsInfinity: false,
					Delay:      1000 * time.Millisecond,
				})
				if err := conn.Push("service/tick", tickdata); err != nil {
					slog.Error("statusTopic err", err)
				}

				time.Sleep(5 * time.Second)

				tickdata, _ = json.Marshal(fluxnode.TickSettings{
					IsInfinity: false,
					Delay:      100 * time.Millisecond,
				})
				if err := conn.Push("service/tick", tickdata); err != nil {
					slog.Error("statusTopic err", err)
				}

				time.Sleep(5 * time.Second)

				tickdata, _ = json.Marshal(fluxnode.TickSettings{
					IsInfinity: true,
					Delay:      100 * time.Millisecond,
				})
				if err := conn.Push("service/tick", tickdata); err != nil {
					slog.Error("statusTopic err", err)
				}

				time.Sleep(5 * time.Second)
			}
		}()

		return nil
	})

	//node.OnReady(func(data settings) error {
	//	logger.Info("ready")
	//	return nil
	//})
	//
	//node.OnStart(func() error {
	//	logger.Info("start")
	//	return nil
	//})
	//
	//node.OnStop(func() error {
	//	logger.Info("stop")
	//	return nil
	//})
	//
	//node.OnRestart(func() error {
	//	logger.Info("restart")
	//	return nil
	//})
	//
	//node.OnError(func(err error) {
	//	logger.Info("error", slog.Any("error", err))
	//})

	node.OnTick(func(d time.Duration, t time.Time) {
		logger.Info("tick", slog.Any("deltaTime", d), slog.Any("timestamp", t))
	})

	if err := node.Run(); err != nil {
		logger.Error("run", slog.Any("error", err))
	}
}
