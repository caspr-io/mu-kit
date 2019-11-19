package rpc

import (
	"fmt"
	"net"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Config struct {
	Port int `split_words:"true" required:"true"`
}

type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
	logger     zerolog.Logger
}

func NewServer(config *Config) (*Server, error) {
	logger := log.Logger.With().Str("component", "µ-kit gRPC").Logger()

	logger.Info().Interface("config", config).Msg("Initializing µ-Kit gRPC server...")

	listener, err := startTcpListener(logger, config)
	if err != nil {
		return nil, err
	}

	return newServer(logger, listener), nil
}

func CreateServer(listener net.Listener) (*Server, error) {
	logger := log.Logger.With().Str("component", "µ-kit-rpc").Logger()

	logger.Info().Msg("Initializing µ-Kit gRPC server...")

	return newServer(logger, listener), nil
}

func newServer(logger zerolog.Logger, listener net.Listener) *Server {
	return &Server{grpcServer: grpc.NewServer(), listener: listener, logger: logger}
}

func (s *Server) Register(service Service) {
	s.grpcServer.RegisterService(service.RPCServiceDesc(), service)
}

func startTcpListener(logger zerolog.Logger, config *Config) (net.Listener, error) {
	address := fmt.Sprintf(":%d", config.Port)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	return listener, nil
}

// Run starts the gRPC server using the configured listener
func (s *Server) Run() {
	s.logger.Info().Msg("Starting µ-Kit gRPC server...")

	if err := s.grpcServer.Serve(s.listener); err != nil {
		s.logger.Fatal().Err(err).Msg("Cannot serve µ-Kit gRPC server")
	}

	s.logger.Info().Msg("µ-Kit gRPC server has shut down")
}

func (s *Server) Close() error {
	s.logger.Info().Msg("Closing µ-Kit gRPC server...")
	s.grpcServer.Stop()

	return nil
}
