package main

import (
	"context"
	"fmt"
	"gossip-glomers/gossip"
	"io"
	"os"
	"path"

	"github.com/google/uuid"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"golang.org/x/exp/slog"
)

const (
	LogDirEnvVarKey  string = "LOG_DIR"
	LogFileExtension string = "txt"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	n := maelstrom.NewNode()
	handler := gossip.New(ctx, n, gossip.Options{Log: getLogger()})

	n.Handle(string(gossip.MessageTypeEcho), handler.HandleEcho)
	n.Handle(string(gossip.MessageTypeGenerate), handler.HandleGenerate)
	n.Handle(string(gossip.MessageTypeTopology), handler.HandleTopology)
	n.Handle(string(gossip.MessageTypeBroadcast), handler.HandleBroadcast)
	n.Handle(string(gossip.MessageTypeAdd), handler.HandleAdd)

	if os.Getenv("GCOUNTER") == "true" {
		n.Handle(string(gossip.MessageTypeRead), handler.HandleReadGCounter)
	} else {
		n.Handle(string(gossip.MessageTypeRead), handler.HandleReadBroadcasts)
	}

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
