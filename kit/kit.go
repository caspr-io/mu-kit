package kit

import (
	"fmt"
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

type MuKitServer struct {
	config  MuServerConfig
	rpc     *rpc.SubSystem
	river   *river.SubSystem
	closeWg sync.WaitGroup
}

// Initizalize the Mu-Kit environment
func initLogger(name string) {
	log.Logger = zerolog.New(os.Stdout).With().Str("service", name).Timestamp().Logger()
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.MessageFieldName = "m"
}

func readConfig(configPrefix string, config interface{}) error {
	try := util.Try()
	try.Try(readConfigFromEnvironment(configPrefix, config))

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

func New(name string, config interface{}) (*MuKitServer, error) {
	initLogger(name)

	if err := readConfig(name, config); err != nil {
		return nil, err
	}

	// Cast to MuKitConfig
	cfg, ok := config.(MuServerConfig)
	if !ok {
		return nil, fmt.Errorf("passed config %T is not a MuKitConfig", config)
	}

	rpcSystem, err := rpc.New(cfg.GrpcConfig())
	if err != nil {
		return nil, err
	}

	riverSystem, err := river.New(cfg.RiverConfig())
	if err != nil {
		return nil, err
	}

	return NewWithSubSystems(cfg, rpcSystem, riverSystem), nil
}

func NewWithSubSystems(config MuServerConfig, rpcSubSystem *rpc.SubSystem, riverSubSystem *river.SubSystem) *MuKitServer {
	return &MuKitServer{config, rpcSubSystem, riverSubSystem, sync.WaitGroup{}}
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

func (s *MuKitServer) Config() MuServerConfig {
	return s.config
}
