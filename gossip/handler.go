package gossip

import (
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"golang.org/x/exp/slog"
)

type Handler struct {
	Node               *maelstrom.Node
	BroadcastMu        *sync.Mutex
	TopologyMu         *sync.Mutex
	Log                *slog.Logger
	topology           map[string][]string
	receivedBroadcasts map[float64]struct{}
}
