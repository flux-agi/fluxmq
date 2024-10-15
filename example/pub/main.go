package main

import (
	"fmt"
	"time"

	"github.com/flux-agi/fluxmq/fluxmq"
)

func main() {
	pub, err := fluxmq.CreatePub()
	if err != nil {
		panic(err)
	}

	for {
		if err := pub.Send("example-topic1", "hello"); err != nil {
			panic(err)
		}
		fmt.Println("send message")
		time.Sleep(250 * time.Millisecond)
	}
}
