package streaming

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog"
)

type MessageContext struct {
	Ctx    context.Context
	Logger *zerolog.Logger
}

// A MessageHandler can optionally implement io.Closer.
// This ensures that its Close() method will be called when the streaming.Subscription is Closed
type MessageHandler interface {
	Name() string
	NewMsg() proto.Message
	Handle(ctx *MessageContext, m proto.Message) error
}
