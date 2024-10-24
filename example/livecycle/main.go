package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
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

	node := fluxnode.Create[settings](ctx, serviceAlias, logger)

	node.OnConnected(func(conn *fluxmq.Connection) error {
		logger.Info("connected")

		go func() {
			for {
				time.Sleep(time.Second)

				readyTopic := fmt.Sprintf("flux.service.%s.ready", serviceAlias)
				dr, _ := json.Marshal(settings{})
				if err := conn.Push(readyTopic, dr); err != nil {
					slog.Error("readyTopic err", err)
				}

				startTopic := fmt.Sprintf("flux.service.%s.start", serviceAlias)
				if err := conn.Push(startTopic, nil); err != nil {
					slog.Error("startTopic err", err)
				}

				stopTopic := fmt.Sprintf("flux.service.%s.stop", serviceAlias)
				if err := conn.Push(stopTopic, nil); err != nil {
					slog.Error("stopTopic err", err)
				}

				restartTopic := fmt.Sprintf("flux.service.%s.restart", serviceAlias)
				if err := conn.Push(restartTopic, nil); err != nil {
					slog.Error("restartTopic err", err)
				}

				errTopic := fmt.Sprintf("flux.service.%s.error", serviceAlias)
				if err := conn.Push(errTopic, nil); err != nil {
					slog.Error("errTopic err", err)
				}
			}
		}()

		return nil
	})

	node.OnReady(func(data settings) error {
		logger.Info("ready")
		return nil
	})

	node.OnStart(func() error {
		logger.Info("start")
		return nil
	})

	node.OnStop(func() error {
		logger.Info("stop")
		return nil
	})

	node.OnRestart(func() error {
		logger.Info("restart")
		return nil
	})

	node.OnError(func(err error) {
		logger.Info("error", slog.Any("error", err))
	})

	node.OnTick(func(d time.Duration, t time.Time) {
		logger.Info("tick", slog.Any("deltaTime", d), slog.Any("timestamp", t))
	})

	if err := node.Run(); err != nil {
		logger.Error("run", slog.Any("error", err))
	}
}
