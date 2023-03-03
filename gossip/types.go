package gossip

type HandlerType string

const (
	HandleTypeEcho      HandlerType = "echo"
	HandleTypeUUID      HandlerType = "generate"
	HandleTypeTopology  HandlerType = "topology"
	HandleTypeRead      HandlerType = "read"
	HandleTypeBroadcast HandlerType = "broadcast"
)
