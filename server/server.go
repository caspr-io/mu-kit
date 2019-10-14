package server

import (
	"fmt"
	"net"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/caspr-io/mu-kit/rpc"
	"google.golang.org/grpc"
)

type MuServerConfig struct {
	GrpcPort int `split_words:"true" required:"true"`
}

type MuServer struct {
	grpcServer *grpc.Server
	config     MuServerConfig
	logger     zerolog.Logger
}

func New() (*MuServer, error) {
	config := MuServerConfig{}
	logger := log.Logger.With().Str("component", "mu-server").Logger()
	err := envconfig.Process("MUKIT", &config)
	if err != nil {
		return nil, err
	}
	logger.Info().Interface("config", config).Send()
	return &MuServer{grpcServer: grpc.NewServer(), config: config, logger: logger}, nil
}

func (s *MuServer) Register(service rpc.Service) {
	s.grpcServer.RegisterService(service.RPCServiceDesc(), service)
}

func (s *MuServer) Run() {
	s.logger.Info().Msg("Starting Mu-Kit Grpc server...")
	address := fmt.Sprintf(":%d", s.config.GrpcPort)
	grpcListener, err := net.Listen("tcp", address)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("Cannot start Mu-Kit listener")
	}
	err = s.grpcServer.Serve(grpcListener)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("Cannot serve Mu-Kit server")
	}
}
