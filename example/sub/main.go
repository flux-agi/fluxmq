package main

import (
	"context"
	"log"

	"github.com/flux-agi/fluxmq/fluxmq"
)

func main() {
	ctx := context.Background()

	conn, err := fluxmq.Connect()
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Subscribe(ctx, "example/pub")
	if err != nil {
		log.Fatal(err)
	}

	for {
		msg := <-ch
		if err := conn.Respond(&msg, []byte("resp")); err != nil {
			log.Fatal(err)
		}
	}
}
