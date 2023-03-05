package gossip

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func (h *Handler) HandleBroadcast(msg maelstrom.Message) error {
	body := messageBody{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	if body.ID == "" {
		body.ID = uuid.New().String()
	}

	if exists := h.saveMessage(body); !exists {
		for _, peer := range h.getPeers() {
			if peer == msg.Src {
				continue
			}
			h.broadcastCh <- broadcastMsg{dst: peer, body: body}
		}
	}

	return h.node.Reply(msg, reply(messageBody{}, MessageTypeBroadcastOK, body.MsgID))
}

func (h *Handler) saveMessage(body messageBody) bool {
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

func (h *Handler) broadcast(ctx context.Context, dst string, body messageBody) {
	ctx, cancel := context.WithTimeout(ctx, h.broadcastTimeout)
	defer cancel()

	_, err := h.node.SyncRPC(ctx, dst, body)
	if err != nil {
		h.log.Info(fmt.Sprintf("%s failed to broadcast to %s", h.ID(), dst))
		h.broadcastCh <- broadcastMsg{dst: dst, body: body}
	}
}
