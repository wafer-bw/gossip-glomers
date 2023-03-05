package gossip

import (
	"encoding/json"

	"github.com/google/uuid"
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func (h Handler) HandleGenerate(msg maelstrom.Message) error {
	body := messageBody{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	return h.node.Reply(msg, reply(messageBody{ID: id.String()}, MessageTypeGenerateOK, body.MsgID))
}
