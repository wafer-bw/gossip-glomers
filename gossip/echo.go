package gossip

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func (h Handler) HandleEcho(msg maelstrom.Message) error {
	body := map[string]any{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	id, ok := body["msg_id"]
	if !ok {
		return ErrMissingMessageID
	}

	body["type"] = MessageTypeEchoOK
	body["in_reply_to"] = id

	return h.Node.Reply(msg, body)
}
