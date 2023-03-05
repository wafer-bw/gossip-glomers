package gossip

import (
	"encoding/json"
)

type MessageType string

const (
	MessageTypeUnknown     MessageType = ""
	MessageTypeEcho        MessageType = "echo"
	MessageTypeEchoOK      MessageType = "echo_ok"
	MessageTypeGenerate    MessageType = "generate"
	MessageTypeGenerateOK  MessageType = "generate_ok"
	MessageTypeTopology    MessageType = "topology"
	MessageTypeRead        MessageType = "read"
	MessageTypeBroadcast   MessageType = "broadcast"
	MessageTypeBroadcastOK MessageType = "broadcast_ok"
)

// Get can read a value from the raw message body into the provided destination returning true if the key was found
// and unmashalling was successful.
func Get[T any](m json.RawMessage, key string, dst T) bool {
	body := map[string]json.RawMessage{}
	if err := json.Unmarshal(m, &body); err != nil {
		return false
	}

	v, ok := body[key]
	if !ok {
		return false
	}

	if err := json.Unmarshal(v, &dst); err != nil {
		return false
	}

	return true
}
