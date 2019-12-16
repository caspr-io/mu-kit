package kit

import (
	"fmt"
	"sync"

	mulog "github.com/caspr-io/mu-kit/log"
	"github.com/caspr-io/mu-kit/rpc"
	"github.com/caspr-io/mu-kit/streaming"
	"github.com/rs/zerolog/log"
)

type MuKitServer struct {
	config    MuServerConfig
	rpcServer *rpc.Server
	river     *streaming.River
	closeWg   sync.WaitGroup
}

func New(name string, config interface{}) (*MuKitServer, error) {
	if err := ReadConfig(name, config); err != nil {
		return nil, err
	}

	// Cast to MuKitConfig
	cfg, ok := config.(MuServerConfig)
	if !ok {
		return nil, fmt.Errorf("passed config %T is not a MuKitConfig", config)
	}

	mulog.Init(name, cfg.LogConfig())

	rpcServer, err := rpc.NewServer(cfg.RPCConfig())
	if err != nil {
		return nil, err
	}

	river, err := streaming.NewRiver(cfg.StreamingConfig())
	if err != nil {
		return nil, err
	}

	return createSystem(cfg, rpcServer, river), nil
}

func createSystem(config MuServerConfig, rpcServer *rpc.Server, river *streaming.River) *MuKitServer {
	return &MuKitServer{config, rpcServer, river, sync.WaitGroup{}}
}

func (s *MuKitServer) Run() {
	// defer s.river.Close()
	s.closeWg.Add(1)

	go s.river.Run()

	go s.rpcServer.Run()

	if err := SignalsHandler(s, log.Logger); err != nil {
		s.Close()
	}

	s.closeWg.Wait()
}

func (s *MuKitServer) Close() error {
	defer s.closeWg.Done()
	s.river.Close()
	s.rpcServer.Close()

	return nil
}

func (s *MuKitServer) River() *streaming.River {
	return s.river
}

func (s *MuKitServer) RPCServer() *rpc.Server {
	return s.rpcServer
}

func (s *MuKitServer) Config() MuServerConfig {
	return s.config
}
