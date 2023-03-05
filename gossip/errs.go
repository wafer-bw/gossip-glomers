package gossip

import "errors"

var (
	ErrMissingMessageID error = errors.New("missing msg_id")
)
