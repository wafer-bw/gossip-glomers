package main

import (
	"context"
	"errors"
	"gossip-glomers/gossip"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

const (
	gCounterKey string = "gcounter"
)

type Handler struct {
	node   *maelstrom.Node
	keyval *maelstrom.KV
}

type message struct {
	maelstrom.MessageBody
	Delta int                `json:"delta,omitempty"`
	Value int                `json:"value,omitempty"`
	Type  gossip.MessageType `json:"type,omitempty"`
}

func New(ctx context.Context, n *maelstrom.Node) *Handler {
	h := &Handler{
		node:   n,
		keyval: maelstrom.NewSeqKV(n),
	}

	go h.initGCounter(ctx)

	return h
}

func (h *Handler) ErrCode(err error) (int, bool) {
	rpcError := &maelstrom.RPCError{}
	if errors.As(err, &rpcError) {
		return rpcError.Code, true
	}
	return 0, false
}
