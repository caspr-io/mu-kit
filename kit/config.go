package kit

import (
	"github.com/caspr-io/mu-kit/river"
	"github.com/caspr-io/mu-kit/rpc"
)

type MuServerConfig interface {
	GrpcConfig() *rpc.SubSystemConfig
	RiverConfig() *river.SubSystemConfig
}

type MuKitConfig struct {
	Grpc  *rpc.SubSystemConfig
	River *river.SubSystemConfig `envconfig:"PUBSUB"`
}

func (c *MuKitConfig) GrpcConfig() *rpc.SubSystemConfig {
	return c.Grpc
}

func (c *MuKitConfig) RiverConfig() *river.SubSystemConfig {
	return c.River
}
