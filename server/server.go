package server

import (
	"fmt"
	"net"

	"github.com/caspr-io/mu-kit/rpc"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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

// New creates a new MuServer
func New() (*MuServer, error) {
	config := MuServerConfig{}
	logger := log.Logger.With().Str("component", "µ-server").Logger()

	if err := envconfig.Process("MUKIT", &config); err != nil {
		return nil, err
	}

	logger.Info().Interface("config", config).Send()

	return &MuServer{grpcServer: grpc.NewServer(), config: config, logger: logger}, nil
}

func (s *MuServer) Register(service rpc.Service) {
	s.grpcServer.RegisterService(service.RPCServiceDesc(), service)
}

// Run starts the gRPC server on the environment provided port number
func (s *MuServer) Run() {
	address := fmt.Sprintf(":%d", s.config.GrpcPort)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("Cannot start µ-Kit listener")
	}

	s.RunWithListener(listener)
}

// RunWithListener starts the gRPC server on the provided listener
func (s *MuServer) RunWithListener(listener net.Listener) {
	s.logger.Info().Msg("Starting µ-Kit gRPC server...")

	if err := s.grpcServer.Serve(listener); err != nil {
		s.logger.Fatal().Err(err).Msg("Cannot serve µ-Kit server")
	}
}
