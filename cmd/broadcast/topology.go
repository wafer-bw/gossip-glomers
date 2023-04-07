package main

import (
	"encoding/json"
	"gossip-glomers/gossip"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func (h *Handler) Topology(msg maelstrom.Message) error {
	body := message{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	h.setTopology(body.Topology)

	reply := struct {
		maelstrom.MessageBody
		Type gossip.MessageType `json:"type"`
	}{
		MessageBody: body.MessageBody,
		Type:        gossip.MessageTypeTopologyOK,
	}

	return h.node.Reply(msg, reply)
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
