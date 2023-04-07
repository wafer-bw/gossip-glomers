package main

import (
	"encoding/json"
	"gossip-glomers/gossip"

	"github.com/google/uuid"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Handler struct {
	node *maelstrom.Node
}

func (h Handler) Generate(msg maelstrom.Message) error {
	body := map[string]any{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	body["type"] = string(gossip.MessageTypeGenerateOK)
	body["id"] = uuid.NewString()

	return h.node.Reply(msg, body)
}
