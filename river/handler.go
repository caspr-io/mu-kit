package river

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog"
)

type MessageContext struct {
	Ctx    context.Context
	Logger *zerolog.Logger
}

// A MuMessageHandler can optionally implement io.Closer.
// This ensures that its Close() method will be called when the river.Subscription is Closed
type MuMessageHandler interface {
	Name() string
	NewMsg() proto.Message
	Handle(ctx *MessageContext, m proto.Message) error
}
