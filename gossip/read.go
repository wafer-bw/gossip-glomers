package gossip

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func (h *Handler) HandleRead(msg maelstrom.Message) error {
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
