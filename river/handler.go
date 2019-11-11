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

type MuMessageHandler interface {
	Name() string
	NewMsg() proto.Message
	Handle(ctx *MessageContext, m proto.Message) error
}
