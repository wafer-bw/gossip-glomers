package gossip

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gossip-glomers/tools/stack"
	"sync"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"golang.org/x/exp/slog"
)

type FaultTolerantBroadcast struct {
	Node               *maelstrom.Node
	BroadcastMu        *sync.Mutex
	TopologyMu         *sync.Mutex
	Log                *slog.Logger
	topology           map[string][]string
	receivedBroadcasts map[float64]struct{}
	unbroadcasted      []broadcastMsg
}

type broadcastMsg struct {
	dst  string
	body map[string]any
}

func (h *FaultTolerantBroadcast) HandleTopology(msg maelstrom.Message) error {
	body := topologyBody{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	h.recordTopology(body.Topology)

	return h.Node.Reply(msg, map[string]any{"type": "topology_ok"})
}

func (h *FaultTolerantBroadcast) HandleRead(msg maelstrom.Message) error {
	body := map[string]any{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	received := h.readBroadcastsReceived()

	return h.Node.Reply(msg, map[string]any{"type": "read_ok", "messages": received})
}

func (h *FaultTolerantBroadcast) HandleBroadcast(msg maelstrom.Message) error {
	body := map[string]any{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	message, ok := body["message"]
	if !ok {
		return errors.New("no message")
	}

	messageNumber, ok := message.(float64)
	if !ok {
		return fmt.Errorf("message is not a float: %s", body["message"])
	}
	h.recordBroadcast(msg.Src, body, messageNumber)

	for _, peer := range h.readTopology()[h.Node.ID()] {
		if peer == msg.Src {
			continue
		}
		broadcast := broadcastMsg{dst: peer, body: body}
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		if _, err := h.Node.SyncRPC(ctx, broadcast.dst, broadcast.body); err != nil {
			h.Log.Info(fmt.Sprintf("broadcast to %s with body %+v failed due to: %s", broadcast.dst, broadcast.body, err))
			h.pushUnbroadcasted(broadcast)
		}
		cancel()
	}

	return h.Node.Reply(msg, map[string]any{"type": "broadcast_ok"})
}

func (h *FaultTolerantBroadcast) recordBroadcast(src string, body map[string]any, value float64) {
	h.BroadcastMu.Lock()
	defer h.BroadcastMu.Unlock()

	if h.receivedBroadcasts == nil {
		h.receivedBroadcasts = map[float64]struct{}{}
	}

	h.receivedBroadcasts[value] = struct{}{}
}

func (h *FaultTolerantBroadcast) readBroadcastsReceived() []float64 {
	h.BroadcastMu.Lock()
	defer h.BroadcastMu.Unlock()

	received := []float64{}
	for k := range h.receivedBroadcasts {
		received = append(received, k)
	}

	return received
}

func (h *FaultTolerantBroadcast) readUnbroadcasted() []broadcastMsg {
	h.BroadcastMu.Lock()
	defer h.BroadcastMu.Unlock()

	return h.unbroadcasted
}

func (h *FaultTolerantBroadcast) recordTopology(topology map[string][]string) {
	h.TopologyMu.Lock()
	defer h.TopologyMu.Unlock()

	h.topology = topology
}

func (h *FaultTolerantBroadcast) readTopology() map[string][]string {
	h.TopologyMu.Lock()
	defer h.TopologyMu.Unlock()

	return h.topology
}

func (h *FaultTolerantBroadcast) broadcastReceived(msg maelstrom.Message, err error) (bool, error) {
	if err != nil {
		return false, err
	}

	body := map[string]any{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return false, err
	}

	if body["type"] == "broadcast_ok" {
		return true, nil
	}

	return false, errors.New("received unknown body type")
}

func (h *FaultTolerantBroadcast) pushUnbroadcasted(broadcast broadcastMsg) {
	h.BroadcastMu.Lock()
	defer h.BroadcastMu.Unlock()

	stack.Push(&h.unbroadcasted, broadcast)
}

func (h *FaultTolerantBroadcast) Blaster(ctx context.Context) {
	// So anyway, I started blasting.
	for _, peer := range h.readTopology()[h.Node.ID()] {
		for _, broadcast := range h.readUnbroadcasted() {
			if ctx.Err() != nil {
				return
			}
			if peer != broadcast.dst {
				continue
			}
			_ = h.Node.Send(broadcast.dst, broadcast.body)
		}
	}
}
