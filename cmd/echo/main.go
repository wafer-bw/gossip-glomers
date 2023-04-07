package main

import (
	"gossip-glomers/gossip"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	handler := Handler{node: n}
	n.Handle(string(gossip.MessageTypeEcho), handler.Echo)

	if err := n.Run(); err != nil {
		panic(err)
	}
}
