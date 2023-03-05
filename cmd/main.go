package main

import (
	"fmt"
	"gossip-glomers/gossip"
	"io"
	"os"
	"path"
	"sync"

	"github.com/google/uuid"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"golang.org/x/exp/slog"
)

const (
	LogDirEnvVarKey  string = "LOG_DIR"
	LogFileExtension string = "txt"
)

func main() {
	n := maelstrom.NewNode()

	handler := &gossip.Handler{
		Node:        n,
		BroadcastMu: &sync.Mutex{},
		TopologyMu:  &sync.Mutex{},
		Log:         getLogger(),
	}

	n.Handle(string(gossip.MessageTypeEcho), handler.HandleEcho)
	n.Handle(string(gossip.MessageTypeGenerate), handler.HandleGenerate)
	n.Handle(string(gossip.MessageTypeTopology), handler.HandleTopology)
	n.Handle(string(gossip.MessageTypeRead), handler.HandleRead)
	n.Handle(string(gossip.MessageTypeBroadcast), handler.HandleBroadcast)

	if err := n.Run(); err != nil {
		panic(err)
	}
}

func getLogger() *slog.Logger {
	dir := os.Getenv(LogDirEnvVarKey)
	if dir == "" {
		return slog.New(slog.NewTextHandler(io.Discard))
	}

	id := uuid.New().String()[9:13]
	fn := fmt.Sprintf("%s.%s", id, LogFileExtension)

	dst, err := os.Create(path.Join(dir, fn))
	if err != nil {
		panic(err)
	}

	return slog.New(slog.NewTextHandler(dst))
}
