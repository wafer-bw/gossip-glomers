package gossip

import (
	"encoding/json"
	"fmt"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"golang.org/x/exp/slog"
)

func (h *Handler) HandleTopology(msg maelstrom.Message) error {
	body := messageBody{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	h.log.Info(
		"topology received",
		slog.String("nodeID", h.node.ID()),
		slog.String("topology", fmt.Sprintf("%+v", body.Topology)),
	)

	h.setTopology(body.Topology)

	return h.node.Reply(msg, reply(messageBody{}, MessageTypeTopologyOK, body.MsgID))
}

func (h *Handler) setTopology(topology map[string][]string) {
	h.topologyMu.Lock()
	defer h.topologyMu.Unlock()

	h.topology = topology
}

func (h *Handler) getPeers() []string {
	h.topologyMu.Lock()
	defer h.topologyMu.Unlock()

	peers, ok := h.topology[h.node.ID()]
	if !ok {
		return []string{}
	}

	return peers
}
