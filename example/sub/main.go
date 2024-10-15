package main

import (
	"context"
	"fmt"

	"github.com/flux-agi/fluxmq/fluxmq"
)

func main() {
	ctx := context.Background()

	sub, err := fluxmq.CreateSub("example_topic")
	if err != nil {
		panic(err)
	}

	dataCh, err := sub.Subscribe(ctx)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case data := <-dataCh:
			fmt.Printf("receive message: %s\n", data)
		}
	}
}
