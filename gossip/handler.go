package gossip

import (
	"context"
	"io"
	"sync"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"golang.org/x/exp/slog"
)

const (
	DefaultBroadcastQueueSize = 1000
	DefaultBroadcastTimeout   = 1 * time.Second
	GCounterKey               = "gcounter"
)

type Options struct {
	Log                *slog.Logger
	BroadcastQueueSize int
	BroadcastTimeout   *time.Duration
}

type Handler struct {
	node             *maelstrom.Node
	messagesMu       *sync.Mutex
	topologyMu       *sync.Mutex
	log              *slog.Logger
	broadcastCh      chan broadcastMsg
	broadcastTimeout time.Duration
	topology         map[string][]string
	messages         map[string]int
	kv               *maelstrom.KV
}

func New(ctx context.Context, node *maelstrom.Node, opts Options) *Handler {
	timeout := DefaultBroadcastTimeout
	if opts.BroadcastTimeout != nil {
		timeout = *opts.BroadcastTimeout
	}

	log := slog.New(slog.NewTextHandler(io.Discard))
	if opts.Log != nil {
		log = opts.Log
	}

	broadcastQueueSize := DefaultBroadcastQueueSize
	if broadcastQueueSize <= 0 {
		broadcastQueueSize = opts.BroadcastQueueSize
	}

	handler := &Handler{
		node:             node,
		messagesMu:       &sync.Mutex{},
		topologyMu:       &sync.Mutex{},
		log:              log,
		broadcastCh:      make(chan broadcastMsg, broadcastQueueSize),
		broadcastTimeout: timeout,
		topology:         map[string][]string{},
		messages:         map[string]int{},
		kv:               maelstrom.NewSeqKV(node),
	}

	go handler.runBroadcaster(ctx)
	go handler.initGCounter(ctx)

	return handler
}

func (h *Handler) ID() string {
	return h.node.ID()
}
