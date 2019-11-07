package kit

import (
	"os"
	"strings"
	"sync"

	"github.com/caspr-io/mu-kit/river"
	"github.com/caspr-io/mu-kit/rpc"
	"github.com/caspr-io/mu-kit/util"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var Config *MuKitConfig

type MuKitServer struct {
	rpc     *rpc.SubSystem
	river   *river.SubSystem
	closeWg sync.WaitGroup
}

type MuKitConfig struct {
	rpcConfig   *rpc.SubSystemConfig
	riverConfig *river.SubSystemConfig
}

// Initizalize the Mu-Kit environment
func initLogger(name string) {
	log.Logger = zerolog.New(os.Stdout).With().Str("service", name).Timestamp().Logger()
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.MessageFieldName = "m"
}

func readConfig(configPrefix string) error {
	try := util.Try()

	rpcConfig := &rpc.SubSystemConfig{}
	riverConfig := &river.SubSystemConfig{}

	try.Try(readConfigFromEnvironment(configPrefix, rpcConfig))
	try.Try(readConfigFromEnvironment(configPrefix, riverConfig))

	Config = &MuKitConfig{rpcConfig: rpcConfig, riverConfig: riverConfig}

	return try.Error()
}

func readConfigFromEnvironment(configPrefix string, config interface{}) func() error {
	return func() error {
		configPrefix = strings.ToUpper(strings.ReplaceAll(configPrefix, "-", "_"))
		if err := envconfig.Process(configPrefix, config); err != nil {
			return err
		}
		return nil
	}
}

func New(name string) (*MuKitServer, error) {
	initLogger(name)

	if err := readConfig(name); err != nil {
		return nil, err
	}

	rpcSystem, err := rpc.New(Config.rpcConfig)
	if err != nil {
		return nil, err
	}

	riverSystem, err := river.New(Config.riverConfig)
	if err != nil {
		return nil, err
	}

	return NewWithSubSystems(rpcSystem, riverSystem), nil
}

func NewWithSubSystems(rpcSubSystem *rpc.SubSystem, riverSubSystem *river.SubSystem) *MuKitServer {
	return &MuKitServer{rpcSubSystem, riverSubSystem, sync.WaitGroup{}}
}

func (s *MuKitServer) Run() {
	// defer s.river.Close()
	s.closeWg.Add(1)

	go s.river.Run()

	go s.rpc.Run()

	if err := SignalsHandler(s, log.Logger); err != nil {
		s.Close()
	}

	s.closeWg.Wait()
}

func (s *MuKitServer) Close() error {
	defer s.closeWg.Done()
	s.river.Close()
	s.rpc.Close()

	return nil
}

func (s *MuKitServer) RiverSystem() *river.SubSystem {
	return s.river
}

func (s *MuKitServer) RPCSystem() *rpc.SubSystem {
	return s.rpc
}
