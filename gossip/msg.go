package gossip

import (
	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

// messageBody is an extension of pkg.go.dev/github.com/jepsen-io/maelstrom/demo/go#messageBody.
type messageBody struct {
	maelstrom.MessageBody
	ID       string              `json:"id,omitempty"`
	Echo     string              `json:"echo,omitempty"`
	Message  int                 `json:"message,omitempty"`
	Messages []int               `json:"messages,omitempty"`
	Value    int                 `json:"value,omitempty"`
	Delta    int                 `json:"delta,omitempty"`
	Topology map[string][]string `json:"topology,omitempty"`
}

type MessageType string

const (
	MessageTypeUnknown     MessageType = ""
	MessageTypeEcho        MessageType = "echo"
	MessageTypeEchoOK      MessageType = "echo_ok"
	MessageTypeGenerate    MessageType = "generate"
	MessageTypeGenerateOK  MessageType = "generate_ok"
	MessageTypeTopology    MessageType = "topology"
	MessageTypeTopologyOK  MessageType = "topology_ok"
	MessageTypeRead        MessageType = "read"
	MessageTypeReadOK      MessageType = "read_ok"
	MessageTypeBroadcast   MessageType = "broadcast"
	MessageTypeBroadcastOK MessageType = "broadcast_ok"
	MessageTypeAdd         MessageType = "add"
	MessageTypeAddOK       MessageType = "add_ok"
)

type broadcastMsg struct {
	dst  string      // destination node
	body messageBody // message body
}

func reply(b messageBody, t MessageType, id int) messageBody {
	b.MsgID = id
	b.InReplyTo = id
	b.Type = string(t)
	return b
}
