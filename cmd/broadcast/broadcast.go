package main

import (
	"context"
	"encoding/json"
	"gossip-glomers/gossip"

	"github.com/google/uuid"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func (h *Handler) Broadcast(msg maelstrom.Message) error {
	body := message{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	if body.ID == "" {
		body.ID = uuid.NewString()
	}

	if exists := h.saveMessage(body); !exists {
		for _, peer := range h.getPeers() {
			if peer == msg.Src {
				continue
			}
			h.broadcastCh <- broadcastMsg{dst: peer, body: body}
		}
	}

	reply := struct {
		maelstrom.MessageBody
		Type gossip.MessageType `json:"type"`
	}{
		MessageBody: body.MessageBody,
		Type:        gossip.MessageTypeBroadcastOK,
	}

	return h.node.Reply(msg, reply)
}

func (h *Handler) saveMessage(body message) bool {
	h.messagesMu.Lock()
	defer h.messagesMu.Unlock()

	if _, ok := h.messages[body.ID]; ok {
		return true
	}

	h.messages[body.ID] = body.Message
	return false
}

func (h *Handler) runBroadcaster(ctx context.Context) {
	for {
		select {
		case broadcast := <-h.broadcastCh:
			h.broadcast(ctx, broadcast.dst, broadcast.body)
		case <-ctx.Done():
			return
		}
	}
}

func (h *Handler) broadcast(ctx context.Context, dst string, body message) {
	ctx, cancel := context.WithTimeout(ctx, h.broadcastTimeout)
	defer cancel()

	_, err := h.node.SyncRPC(ctx, dst, body)
	if err != nil {
		h.broadcastCh <- broadcastMsg{dst: dst, body: body}
	}
}
