package gossip

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"golang.org/x/exp/slog"
)

type topologyBody struct {
	Type     string              `json:"type"`
	Topology map[string][]string `json:"topology"`
}

type broadcastMsg struct {
	uid  string
	dst  string
	body map[string]any
}

func (h *Handler) HandleTopology(msg maelstrom.Message) error {
	body := topologyBody{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	h.recordTopology(body.Topology)

	h.Log.Info(
		"read topology",
		slog.String("id", h.Node.ID()),
		slog.String("topo", fmt.Sprintf("%+v", body.Topology)),
	)

	return h.Node.Reply(msg, map[string]any{"type": "topology_ok"})
}

func (h *Handler) HandleRead(msg maelstrom.Message) error {
	body := map[string]any{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	received := h.readBroadcastsReceived()

	return h.Node.Reply(msg, map[string]any{"type": "read_ok", "messages": received})
}

func (h *Handler) HandleBroadcast(msg maelstrom.Message) error {
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

	if alreadyReceived := h.recordBroadcast(msg.Src, messageNumber); alreadyReceived {
		return h.Node.Reply(msg, map[string]any{"type": "broadcast_ok"})
	}

	h.doBroadcasting(msg, body)

	return h.Node.Reply(msg, map[string]any{"type": "broadcast_ok"})
}

func (h *Handler) recordBroadcast(src string, value float64) bool {
	h.BroadcastMu.Lock()
	defer h.BroadcastMu.Unlock()

	if h.receivedBroadcasts == nil {
		h.receivedBroadcasts = map[float64]struct{}{}
	}

	if _, ok := h.receivedBroadcasts[value]; ok {
		return true
	}

	h.receivedBroadcasts[value] = struct{}{}
	return false
}

func (h *Handler) readBroadcastsReceived() []float64 {
	h.BroadcastMu.Lock()
	defer h.BroadcastMu.Unlock()

	received := []float64{}
	for k := range h.receivedBroadcasts {
		received = append(received, k)
	}

	return received
}

func (h *Handler) recordTopology(topology map[string][]string) {
	h.TopologyMu.Lock()
	defer h.TopologyMu.Unlock()

	h.topology = topology
}

func (h *Handler) readTopology() map[string][]string {
	h.TopologyMu.Lock()
	defer h.TopologyMu.Unlock()

	return h.topology
}

func (h *Handler) doBroadcasting(msg maelstrom.Message, body map[string]any) {
	uid := uuid.New().String()
	for _, peer := range h.readTopology()[h.Node.ID()] {
		if peer == msg.Src {
			continue
		}
		broadcast := broadcastMsg{dst: peer, uid: uid, body: body}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		_, err := h.Node.SyncRPC(ctx, broadcast.dst, broadcast.body)
		if err != nil {
			h.Log.Info(fmt.Sprintf("%s failed to broadcast to %s", h.Node.ID(), broadcast.dst))
			go h.rebroadcast(broadcast)
		}
		cancel()
	}
}

func (h *Handler) rebroadcast(broadcast broadcastMsg) {
	retry := 0
	maxRetry := 100
	for {
		time.Sleep(500 * time.Millisecond)
		retry++
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		_, err := h.Node.SyncRPC(ctx, broadcast.dst, broadcast.body)
		cancel()
		if err == nil {
			h.Log.Info(fmt.Sprintf("%s rebroadcasted %+v to %s after %d attempts", h.Node.ID(), broadcast.body, broadcast.dst, retry))
			return
		}
		if retry > maxRetry {
			h.Log.Info(fmt.Sprintf("%s failed to rebroadcast %+v to %s after %d attempts", h.Node.ID(), broadcast.body, broadcast.dst, retry))
			return
		}
	}
}
