package main

import (
	"fmt"
	"gossip-glomers/gossip"
	"io"
	"os"
	"sync"

	"github.com/google/uuid"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"golang.org/x/exp/slog"
)

func main() {
	n := maelstrom.NewNode()

	broadcast := gossip.Broadcast{
		BroadcastMu: &sync.Mutex{},
		TopologyMu:  &sync.Mutex{},
		Node:        n,
		Log:         getLogger(),
	}

	n.Handle(string(gossip.HandleTypeEcho), gossip.Echo{Node: n}.Handle)
	n.Handle(string(gossip.HandleTypeUUID), gossip.UUID{Node: n}.Handle)
	n.Handle(string(gossip.HandleTypeTopology), broadcast.HandleTopology)
	n.Handle(string(gossip.HandleTypeRead), broadcast.HandleRead)
	n.Handle(string(gossip.HandleTypeBroadcast), broadcast.HandleBroadcast)

	if err := n.Run(); err != nil {
		panic(err)
	}
}

func getLogger() *slog.Logger {
	dir := os.Getenv("LOG_DIR")
	if dir == "" {
		return slog.New(slog.NewTextHandler(io.Discard))
	}

	id := uuid.New().String()[9:13]
	fn := fmt.Sprintf("%s.txt", id)

	dst, err := os.Create(fmt.Sprintf("%s/%s", dir, fn))
	if err != nil {
		panic(err)
	}

	return slog.New(slog.NewTextHandler(dst))
}
