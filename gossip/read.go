package gossip

import (
	"context"
	"encoding/json"
	"fmt"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func (h *Handler) HandleReadBroadcasts(msg maelstrom.Message) error {
	body := messageBody{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	msgs := h.getMessages()

	return h.node.Reply(msg, reply(messageBody{Messages: msgs}, MessageTypeReadOK, body.MsgID))
}

func (h *Handler) getMessages() []int {
	h.messagesMu.Lock()
	defer h.messagesMu.Unlock()

	received := []int{}
	for _, v := range h.messages {
		received = append(received, v)
	}

	return received
}

func (h *Handler) HandleReadGCounter(msg maelstrom.Message) error {
	ctx := context.Background()

	body := messageBody{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	v, err := h.kv.ReadInt(ctx, GCounterKey)
	if err != nil {
		return err
	}

	h.log.Info(fmt.Sprintf("read %d", v))

	return h.node.Reply(msg, reply(messageBody{Value: v}, MessageTypeReadOK, body.MsgID))
}
