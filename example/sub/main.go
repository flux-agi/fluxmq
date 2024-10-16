package main

import (
	"context"
	"fmt"
	"log"

	"github.com/flux-agi/fluxmq/fluxmq"
)

func main() {
	ctx := context.Background()

	conn, err := fluxmq.Connect()
	if err != nil {
		log.Fatal(err)
	}

	fef, err := conn.CreateSub("example/pub")
	ch, err := fef.Recv(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for {
		msg := <-ch
		fmt.Printf("recv message: %s\n", msg.Payload)
	}
}
