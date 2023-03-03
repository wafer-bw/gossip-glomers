package main

import (
	"gossip-glomers/gossip"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	broadcast := gossip.Broadcast{
		BroadcastMu:   &sync.Mutex{},
		TopologyMu:    &sync.Mutex{},
		Propogate:     true,
		Partitionable: true,
		Node:          n,
	}

	n.Handle(string(gossip.HandleTypeTopology), broadcast.HandleTopology)
	n.Handle(string(gossip.HandleTypeRead), broadcast.HandleRead)
	n.Handle(string(gossip.HandleTypeBroadcast), broadcast.HandleBroadcast)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
