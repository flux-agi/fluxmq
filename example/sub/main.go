package main

import (
	"context"
	"fmt"

	"zenoh_node_1"
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
