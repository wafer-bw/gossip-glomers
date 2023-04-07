package main

import (
	"context"
	"encoding/json"
	"gossip-glomers/gossip"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func (h *Handler) Add(msg maelstrom.Message) error {
	ctx := context.Background()

	body := message{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	if body.Delta == 0 {
		reply := struct {
			maelstrom.MessageBody
			Type gossip.MessageType `json:"type"`
		}{
			MessageBody: body.MessageBody,
			Type:        gossip.MessageTypeAddOK,
		}
		return h.node.Reply(msg, reply)
	}

	// CAS until we succeed
	for {
		value, err := h.keyval.ReadInt(ctx, gCounterKey)
		if err != nil {
			return err
		}

		if err := h.keyval.CompareAndSwap(ctx, gCounterKey, value, value+body.Delta, false); err == nil {
			break
		}
	}

	reply := struct {
		maelstrom.MessageBody
		Type gossip.MessageType `json:"type"`
	}{
		MessageBody: body.MessageBody,
		Type:        gossip.MessageTypeAddOK,
	}

	return h.node.Reply(msg, reply)
}

func (h *Handler) initGCounter(ctx context.Context) {
	// Wait for node ID to be assigned
	for h.node.ID() == "" {
		time.Sleep(10 * time.Millisecond)
	}

	// Only n0 initializes the gcounter
	if h.node.ID() != "n0" {
		return
	}

	for {
		_, err := h.keyval.Read(ctx, gCounterKey)
		if c, ok := h.ErrCode(err); err != nil && ok && c == maelstrom.KeyDoesNotExist {
			if err := h.keyval.Write(ctx, gCounterKey, 0); err != nil {
				continue
			}
			return
		}
	}
}
