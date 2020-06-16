package streaming

import (
	"context"

	"google.golang.org/protobuf/proto"
)

// A MessageHandler can optionally implement io.Closer.
// This ensures that its Close() method will be called when the streaming.Subscription is Closed
type MessageHandler interface {
	Name() string
	NewMsg() proto.Message
	Handle(ctx context.Context, m proto.Message) error
}
