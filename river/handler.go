package river

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/rs/zerolog"
)

type MessageContext struct {
	ctx    context.Context
	logger *zerolog.Logger
}

type MuMessageHandler interface {
	Name() string
	NewMsg() proto.Message
	Handle(ctx *MessageContext, m proto.Message) error
}
