package main

import (
	"log"
	"time"

	"github.com/flux-agi/fluxmq/pkg"
)

func main() {
	conn, err := pkg.Connect()
	if err != nil {
		log.Fatal(err)
	}

	fef, err := conn.CreatePub()
	for {
		if err := fef.Push("example/pub", []byte("big bo")); err != nil {
			log.Fatal(err)
		}

		time.Sleep(250 * time.Millisecond)
	}
}
