package main

import (
	"gossip-glomers/gossip"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	n.Handle(string(gossip.HandleTypeEcho), gossip.Echo{Node: n}.Handle)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
