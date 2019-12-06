package kit

import (
	"strings"

	"github.com/caspr-io/mu-kit/rpc"
	"github.com/caspr-io/mu-kit/streaming"
	"github.com/kelseyhightower/envconfig"
)

type MuServerConfig interface {
	RPCConfig() *rpc.Config
	StreamingConfig() *streaming.Config
}

type MuKitConfig struct {
	Grpc      *rpc.Config
	Streaming *streaming.Config
}

func (c *MuKitConfig) RPCConfig() *rpc.Config {
	return c.Grpc
}

func (c *MuKitConfig) StreamingConfig() *streaming.Config {
	return c.Streaming
}

func ReadConfig(configPrefix string, config interface{}) error {
	configPrefix = strings.ToUpper(strings.ReplaceAll(configPrefix, "-", "_"))
	if err := envconfig.Process(configPrefix, config); err != nil {
		return err
	}

	return nil
}
