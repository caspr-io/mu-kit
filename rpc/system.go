package rpc

import (
	"fmt"
	"net"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type SubSystemConfig struct {
	GrpcPort int `split_words:"true" required:"true"`
}

type SubSystem struct {
	grpcServer *grpc.Server
	listener   net.Listener
	logger     zerolog.Logger
}

// New creates a new RPC SubSystem
func New(config *SubSystemConfig) (*SubSystem, error) {
	logger := log.Logger.With().Str("component", "µ-kit-rpc").Logger()

	logger.Info().Interface("config", config).Msg("Initializing µ-Kit RPC subsystem...")

	listener, err := startTcpListener(logger, config)
	if err != nil {
		return nil, err
	}

	return newServer(logger, listener), nil
}

func NewWithListener(listener net.Listener) (*SubSystem, error) {
	logger := log.Logger.With().Str("component", "µ-kit-rpc").Logger()

	logger.Info().Msg("Initializing µ-Kit RPC subsystem...")

	return newServer(logger, listener), nil
}

func newServer(logger zerolog.Logger, listener net.Listener) *SubSystem {
	return &SubSystem{grpcServer: grpc.NewServer(), listener: listener, logger: logger}
}

func (s *SubSystem) Register(service Service) {
	s.grpcServer.RegisterService(service.RPCServiceDesc(), service)
}

func startTcpListener(logger zerolog.Logger, config *SubSystemConfig) (net.Listener, error) {
	address := fmt.Sprintf(":%d", config.GrpcPort)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	return listener, nil
}

// Run starts the gRPC server using the configured listener
func (s *SubSystem) Run() {
	s.logger.Info().Msg("Starting µ-Kit gRPC server...")

	if err := s.grpcServer.Serve(s.listener); err != nil {
		s.logger.Fatal().Err(err).Msg("Cannot serve µ-Kit server")
	}
}
