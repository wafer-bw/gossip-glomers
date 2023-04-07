package gossip

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
