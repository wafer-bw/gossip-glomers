package main

import (
	"encoding/json"
	"gossip-glomers/gossip"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func (h *Handler) Read(msg maelstrom.Message) error {
	body := message{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	reply := struct {
		maelstrom.MessageBody
		Messages []int              `json:"messages"`
		Type     gossip.MessageType `json:"type"`
	}{
		MessageBody: body.MessageBody,
		Messages:    h.getMessages(),
		Type:        gossip.MessageTypeReadOK,
	}

	return h.node.Reply(msg, reply)
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
