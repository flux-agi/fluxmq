package main

import (
	"context"
	"log"

	"github.com/flux-agi/fluxmq/pkg"
)

func main() {
	ctx := context.Background()

	conn, err := pkg.Connect()
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
		if err := conn.Respond(&msg, []byte("resp")); err != nil {
			log.Fatal(err)
		}
	}
}
