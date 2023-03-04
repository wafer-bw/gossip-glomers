package gossip

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Broadcast struct {
	Node               *maelstrom.Node
	BroadcastMu        *sync.Mutex
	TopologyMu         *sync.Mutex
	Propogate          bool
	Partitionable      bool
	topology           map[string][]string
	receivedBroadcasts map[float64]struct{}
}

type topologyBody struct {
	Type     string              `json:"type"`
	Topology map[string][]string `json:"topology"`
}

func (h *Broadcast) HandleTopology(msg maelstrom.Message) error {
	body := topologyBody{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	h.recordTopology(body.Topology)

	return h.Node.Reply(msg, map[string]any{"type": "topology_ok"})
}

func (h *Broadcast) HandleRead(msg maelstrom.Message) error {
	body := map[string]any{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	received := h.readBroadcasts()

	return h.Node.Reply(msg, map[string]any{"type": "read_ok", "messages": received})
}

func (h *Broadcast) HandleBroadcast(msg maelstrom.Message) error {
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

	alreadyReceived := h.recordBroadcast(messageNumber)
	if alreadyReceived {
		return h.Node.Reply(msg, map[string]any{"type": "broadcast_ok"})
	}

	if h.Propogate {
		topology := h.readTopology()
		for _, peer := range topology[h.Node.ID()] {
			if peer == msg.Src {
				continue
			}

			_ = h.Node.Send(peer, body)
		}
	}

	return h.Node.Reply(msg, map[string]any{"type": "broadcast_ok"})
}

func (h *Broadcast) recordBroadcast(broadcast float64) bool {
	h.BroadcastMu.Lock()
	defer h.BroadcastMu.Unlock()

	if h.receivedBroadcasts == nil {
		h.receivedBroadcasts = map[float64]struct{}{}
	}

	if _, ok := h.receivedBroadcasts[broadcast]; ok {
		return true
	}

	h.receivedBroadcasts[broadcast] = struct{}{}
	return false
}

func (h *Broadcast) readBroadcasts() []float64 {
	h.BroadcastMu.Lock()
	defer h.BroadcastMu.Unlock()

	received := []float64{}
	for k := range h.receivedBroadcasts {
		received = append(received, k)
	}

	return received
}

func (h *Broadcast) recordTopology(topology map[string][]string) {
	h.TopologyMu.Lock()
	defer h.TopologyMu.Unlock()

	h.topology = topology
}

func (h *Broadcast) readTopology() map[string][]string {
	h.TopologyMu.Lock()
	defer h.TopologyMu.Unlock()

	return h.topology
}
