package gossip

import (
	"encoding/json"
	"fmt"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

type Echo struct {
	Node *maelstrom.Node
}

func (h Echo) Handle(msg maelstrom.Message) error {
	body := map[string]any{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	body["type"] = "echo_ok"

	if replyTo, ok := body["msg_id"]; ok {
		body["in_reply_to"] = replyTo
	} else {
		return fmt.Errorf("missing msg_id in body: %s", body)
	}

	return h.Node.Reply(msg, body)
}
