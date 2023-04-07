package main

import (
	"context"
	"gossip-glomers/gossip"
	"sync"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

const (
	broadcastQueueSize int           = 100
	broadcastTimeout   time.Duration = 1 * time.Second
)

type Handler struct {
	node             *maelstrom.Node
	messagesMu       *sync.Mutex
	topologyMu       *sync.Mutex
	broadcastCh      chan broadcastMsg
	broadcastTimeout time.Duration
	topology         map[string][]string
	messages         map[string]int
}

type message struct {
	maelstrom.MessageBody
	ID       string              `json:"id,omitempty"`
	Message  int                 `json:"message,omitempty"`
	Messages []int               `json:"messages,omitempty"`
	Topology map[string][]string `json:"topology,omitempty"`
	Type     gossip.MessageType  `json:"type,omitempty"`
}

type broadcastMsg struct {
	dst  string // destination node
	body message
}

func New(ctx context.Context, n *maelstrom.Node) *Handler {
	h := &Handler{
		node:             n,
		messagesMu:       &sync.Mutex{},
		topologyMu:       &sync.Mutex{},
		broadcastCh:      make(chan broadcastMsg, broadcastQueueSize),
		broadcastTimeout: broadcastTimeout,
		topology:         map[string][]string{},
		messages:         map[string]int{},
	}

	go h.runBroadcaster(ctx)

	return h
}
