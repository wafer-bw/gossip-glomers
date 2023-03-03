package gossip

import (
	"encoding/json"

	"github.com/google/uuid"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type UUID struct {
	Node *maelstrom.Node
}

func (h UUID) Handle(msg maelstrom.Message) error {
	body := map[string]any{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	body["type"] = "generate_ok"
	body["id"] = id

	return h.Node.Reply(msg, body)
}
