package main

import (
	"context"
	"gossip-glomers/gossip"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	n := maelstrom.NewNode()

	handler := New(ctx, n)
	n.Handle(string(gossip.MessageTypeAdd), handler.Add)
	n.Handle(string(gossip.MessageTypeRead), handler.Read)

	if err := n.Run(); err != nil {
		panic(err)
	}
}
