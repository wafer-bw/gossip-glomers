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

	// TODO: For some reason this doesn't heal in time all the time
	for {
		currentValue, _ := h.keyval.ReadInt(ctx, gCounterKey)
		newValue := currentValue + body.Delta
		err := h.keyval.CompareAndSwap(ctx, gCounterKey, currentValue, newValue, false)
		if c, ok := h.ErrCode(err); err != nil && ok && c == maelstrom.PreconditionFailed {
			continue
		} else if err != nil {
			return err
		}

		break
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
