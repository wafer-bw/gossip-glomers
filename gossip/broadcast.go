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
	Mutex              *sync.Mutex
	receivedBroadcasts []float64
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

	nodeIDs := h.Node.NodeIDs()
	for _, nodeID := range nodeIDs {
		if _, ok := body.Topology[nodeID]; !ok {
			return errors.New("topology is missing node " + nodeID)
		}
	}

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

	h.recordBroadcast(messageNumber)

	return h.Node.Reply(msg, map[string]any{"type": "broadcast_ok"})
}

func (h *Broadcast) recordBroadcast(broadcast float64) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	h.receivedBroadcasts = append(h.receivedBroadcasts, broadcast)
}

func (h *Broadcast) readBroadcasts() []float64 {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	return h.receivedBroadcasts
}
