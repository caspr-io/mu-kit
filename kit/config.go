package kit

import (
	"github.com/caspr-io/mu-kit/rpc"
	"github.com/caspr-io/mu-kit/streaming"
)

type MuServerConfig interface {
	RpcConfig() *rpc.Config
	StreamingConfig() *streaming.Config
}

type MuKitConfig struct {
	Grpc      *rpc.Config
	Streaming *streaming.Config
}

func (c *MuKitConfig) RpcConfig() *rpc.Config {
	return c.Grpc
}

func (c *MuKitConfig) StreamingConfig() *streaming.Config {
	return c.Streaming
}