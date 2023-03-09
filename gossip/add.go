package gossip

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	"golang.org/x/exp/slog"
)

func (h *Handler) HandleAdd(msg maelstrom.Message) error {
	ctx := context.Background()
	s := time.Now()

	body := messageBody{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	if body.Delta == 0 {
		return h.node.Reply(msg, reply(messageBody{}, MessageTypeAddOK, body.MsgID))
	}

	// TODO: For some reason this doesn't heal in time all the time
	attempt := 0
	for {
		attempt++
		v, _ := h.kv.ReadInt(ctx, GCounterKey)
		h.log.Info(fmt.Sprintf("attempting (%d) to add %d to %d", attempt, body.Delta, v))
		err := h.kv.CompareAndSwap(ctx, GCounterKey, v, v+body.Delta, false)
		if c, ok := h.ErrCode(err); err != nil && ok && c == maelstrom.PreconditionFailed {
			continue
		} else if err != nil {
			return err
		} else {
			h.log.Info(fmt.Sprintf("added %d to %d = %d", body.Delta, v, v+body.Delta), slog.Duration("elapsed", time.Since(s)))
			break
		}
	}

	return h.node.Reply(msg, reply(messageBody{}, MessageTypeAddOK, body.MsgID))
}

func (h *Handler) ErrCode(err error) (int, bool) {
	rpcError := &maelstrom.RPCError{}
	if errors.As(err, &rpcError) {
		return rpcError.Code, true
	}
	return 0, false
}

func (h *Handler) initGCounter(ctx context.Context) {
	for h.ID() == "" {
		h.log.Info("waiting to be assigned a node ID")
	}
	if h.ID() != "n0" {
		h.log.Info("not n0, not initializing gcounter", slog.String("node", h.ID()))
		return
	}
	h.log.Info("initializing gcounter", slog.String("node", h.ID()))

	for {
		// TODO: Any node could do this and they could overwrite eachother.
		_, err := h.kv.Read(ctx, GCounterKey)
		if c, ok := h.ErrCode(err); err != nil && ok && c == maelstrom.KeyDoesNotExist {
			if err := h.kv.Write(ctx, GCounterKey, 0); err != nil {
				continue
			}
			h.log.Info("initialized gcounter")
			return
		}
	}
}
