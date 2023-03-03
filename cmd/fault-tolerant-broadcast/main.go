package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"gossip-glomers/gossip"

	"github.com/google/uuid"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"golang.org/x/exp/slog"
)

func main() {
	n := maelstrom.NewNode()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// TODO: make this path come from an env var
	logDst, err := os.Create(fmt.Sprintf("/Users/bwakeford/Documents/dev/gossip-glomers/log/%s-log.txt", uuid.New().String()))
	if err != nil {
		panic(err)
	}

	log := slog.New(slog.NewTextHandler(logDst))

	broadcast := gossip.FaultTolerantBroadcast{
		BroadcastMu: &sync.Mutex{},
		TopologyMu:  &sync.Mutex{},
		Node:        n,
		Log:         log,
	}

	n.Handle(string(gossip.HandleTypeTopology), broadcast.HandleTopology)
	n.Handle(string(gossip.HandleTypeRead), broadcast.HandleRead)
	n.Handle(string(gossip.HandleTypeBroadcast), broadcast.HandleBroadcast)

	go broadcast.Blaster(ctx)

	if err := n.Run(); err != nil {
		panic(err)
	}
}
