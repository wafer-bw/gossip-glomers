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

	// The maelstrom SeqKV store isn't actually sequential.
	// It likely uses a different lock for reads vs writes.
	// To get around this we do a CAS with the value we read.
	// If the CAS succeeds we know we have the latest value.
	// If the CAS fails we know we read a stale value.
	if err := h.keyval.CompareAndSwap(ctx, gCounterKey, value, value, false); err != nil {
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
