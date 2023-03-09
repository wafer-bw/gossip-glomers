package gossip

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func (h *Handler) HandleAdd(msg maelstrom.Message) error {
	ctx := context.Background()

	body := messageBody{}
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	if body.Delta == 0 {
		return h.node.Reply(msg, reply(messageBody{}, MessageTypeAddOK, body.MsgID))
	}

	// TODO: Any node could do this and they could overwrite eachother.
	_, err := h.kv.Read(ctx, GCounterKey)
	if c, ok := h.ErrCode(err); ok && c == maelstrom.KeyDoesNotExist {
		if err := h.kv.Write(ctx, GCounterKey, 0); err != nil {
			return err
		}
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
			h.log.Info(fmt.Sprintf("added %d to %d = %d", body.Delta, v, v+body.Delta))
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
