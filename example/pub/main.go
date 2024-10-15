package main

import (
	"fmt"
	"time"

	"zenoh_node_1"
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
