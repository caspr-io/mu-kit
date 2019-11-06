package kit

import (
	"os"
	"sync"

	"github.com/caspr-io/mu-kit/river"
	"github.com/caspr-io/mu-kit/rpc"
	"github.com/caspr-io/mu-kit/util"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var onceInit sync.Once
var onceConfig sync.Once
var Config *MuKitConfig

type MuKitServer struct {
	rpc   *rpc.SubSystem
	river *river.SubSystem
}

type MuKitConfig struct {
	rpcConfig   *rpc.SubSystemConfig
	riverConfig *river.SubSystemConfig
}

// Initizalize the Mu-Kit environment
func initLogger() {
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"
}

func readConfig(configPrefix string) error {
	try := util.Try()

	onceConfig.Do(func() {
		rpcConfig := &rpc.SubSystemConfig{}
		riverConfig := &river.SubSystemConfig{}
		try.Try(readConfigFromEnvironment(configPrefix, rpcConfig))
		try.Try(readConfigFromEnvironment(configPrefix, riverConfig))
		Config = &MuKitConfig{rpcConfig: rpcConfig, riverConfig: riverConfig}
	})

	return try.Error()
}

func readConfigFromEnvironment(configPrefix string, config interface{}) func() error {
	return func() error {
		if err := envconfig.Process(configPrefix, config); err != nil {
			return err
		}
		return nil
	}
}

func New(configPrefix string) (*MuKitServer, error) {
	initLogger()

	if err := readConfig(configPrefix); err != nil {
		return nil, err
	}

	rpcSystem, err := rpc.New(Config.rpcConfig)
	if err != nil {
		return nil, err
	}

	riverSystem, err := river.New(&log.Logger, Config.riverConfig)
	if err != nil {
		return nil, err
	}

	return NewWithSubSystems(rpcSystem, riverSystem), nil
}

func NewWithSubSystems(rpcSubSystem *rpc.SubSystem, riverSubSystem *river.SubSystem) *MuKitServer {
	return &MuKitServer{rpcSubSystem, riverSubSystem}
}

func (s *MuKitServer) Run() {
	defer s.river.Close()

	go s.river.Run()

	s.rpc.Run()
}

func (s *MuKitServer) RiverSystem() *river.SubSystem {
	return s.river
}
