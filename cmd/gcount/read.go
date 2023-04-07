package main

import (
	"context"
	"encoding/json"
	"gossip-glomers/gossip"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func (h *Handler) Read(msg maelstrom.Message) error {
	ctx := context.Background()

	body := message{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	value, err := h.keyval.ReadInt(ctx, gCounterKey)
	if err != nil {
		return err
	}

	reply := struct {
		maelstrom.MessageBody
		Value int                `json:"value"`
		Type  gossip.MessageType `json:"type"`
	}{
		MessageBody: body.MessageBody,
		Value:       value,
		Type:        gossip.MessageTypeReadOK,
	}

	return h.node.Reply(msg, reply)
}
