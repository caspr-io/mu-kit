package river

import (
	"context"

	"github.com/golang/protobuf/proto"
)

type MuMessageHandler interface {
	Name() string
	NewMsg() proto.Message
	Handle(ctx context.Context, m proto.Message) error
}
