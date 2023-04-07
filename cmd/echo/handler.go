package main

import (
	"encoding/json"
	"gossip-glomers/gossip"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Handler struct {
	node *maelstrom.Node
}

func (h Handler) Echo(msg maelstrom.Message) error {
	body := map[string]any{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	body["type"] = string(gossip.MessageTypeEchoOK)

	return h.node.Reply(msg, body)
}
