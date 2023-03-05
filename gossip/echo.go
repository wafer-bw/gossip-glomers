package gossip

import (
	"encoding/json"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func (h Handler) HandleEcho(msg maelstrom.Message) error {
	body := messageBody{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	return h.node.Reply(msg, reply(messageBody{Echo: body.Echo}, MessageTypeEchoOK, body.MsgID))
}
