package singlenode

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func Run() {
	n := maelstrom.NewNode()

	app := &app{
		node:               n,
		receivedBroadcasts: []float64{},
		mutex:              &sync.Mutex{},
	}

	n.Handle("broadcast", app.broadcastHandler)
	n.Handle("read", app.readHandler)
	n.Handle("topology", app.topologyHandler)

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

type app struct {
	node               *maelstrom.Node
	mutex              *sync.Mutex
	receivedBroadcasts []float64
}

func (a *app) recordBroadcast(broadcast float64) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.receivedBroadcasts = append(a.receivedBroadcasts, broadcast)
}

func (a *app) ReceivedBroadcasts() []float64 {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.receivedBroadcasts
}

func (a *app) broadcastHandler(msg maelstrom.Message) error {
	body := map[string]any{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	message, ok := body["message"]
	if !ok {
		return errors.New("no message")
	}

	messageNumber, ok := message.(float64)
	if !ok {
		return fmt.Errorf("message is not a float: %s", body["message"])
	}

	a.recordBroadcast(messageNumber)

	return a.node.Reply(msg, map[string]any{"type": "broadcast_ok"})
}

func (a *app) readHandler(msg maelstrom.Message) error {
	body := map[string]any{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	received := a.ReceivedBroadcasts()

	return a.node.Reply(msg, map[string]any{"type": "read_ok", "messages": received})
}

func (a *app) topologyHandler(msg maelstrom.Message) error {
	body := topologyBody{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	nodeIDs := a.node.NodeIDs()
	for _, nodeID := range nodeIDs {
		if _, ok := body.Topology[nodeID]; !ok {
			return errors.New("topology is missing node " + nodeID)
		}
	}

	return a.node.Reply(msg, map[string]any{"type": "topology_ok"})
}

type topologyBody struct {
	Type     string              `json:"type"`
	Topology map[string][]string `json:"topology"`
}
