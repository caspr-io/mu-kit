package rpc

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type Config struct {
	Port    int `split_words:"true" required:"true"`
	WebPort int `split_words:"true" required:"false"`
}

type Server struct {
	grpcAddress    string
	grpcWebAddress string
	grpcServer     *grpc.Server
	grpcListener   net.Listener
	grpcWebServer  *http.Server
	logger         zerolog.Logger
}

// Creates a gRPC Server with an optional gRPC-Web proxy
func NewServer(config *Config) (*Server, error) {
	logger := log.Logger.With().Str("component", "rpc").Logger()

	logger.Info().Interface("config", config).Msg("Initializing gRPC server...")

	grpcServer := grpc.NewServer()

	grpcAddress := fmt.Sprintf(":%d", config.Port)

	listener, err := newGrpcListener(grpcAddress)
	if err != nil {
		return nil, err
	}

	var grpcWebServer *http.Server

	var grpcWebAddress string

	if config.WebPort != 0 {
		grpcWebAddress = fmt.Sprintf(":%d", config.WebPort)
		grpcWebServer = newGrpcWebServer(grpcServer, grpcWebAddress)
	} else {
		grpcWebServer = nil
	}

	return &Server{
			grpcServer:     grpcServer,
			grpcWebAddress: grpcWebAddress,
			grpcAddress:    grpcAddress,
			grpcListener:   listener,
			grpcWebServer:  grpcWebServer,
			logger:         logger},
		nil
}

// Creates a gRPC server for the supplied listener. A gRPC-Web proxy is not created
func NewTestServer(listener net.Listener) (*Server, error) {
	logger := log.Logger.With().Str("component", "rpc").Logger()

	logger.Info().Msg("Initializing test gRPC server...")

	return &Server{
		grpcAddress:    "localhost",
		grpcWebAddress: "",
		grpcServer:     grpc.NewServer(),
		grpcListener:   listener,
		grpcWebServer:  nil,
		logger:         logger}, nil
}

func newGrpcListener(address string) (net.Listener, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}

	return listener, nil
}

func newGrpcWebServer(grpcServer *grpc.Server, address string) *http.Server {
	wrappedGrpc := grpcweb.WrapServer(grpcServer,
		grpcweb.WithCorsForRegisteredEndpointsOnly(false), // TODO: Figure this out for tighter security
		grpcweb.WithOriginFunc(func(origin string) bool { // TODO: Figure this out for tighter security
			return true
		}))

	grpcWebHandler := http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		logger := log.With().
			Str("url", req.URL.String()).
			Str("method", req.Method).
			Str("content-type", req.Header.Get("content-type")).
			Logger()

		if wrappedGrpc.IsAcceptableGrpcCorsRequest(req) {
			logger.Debug().Msg("is a gRPC-Web CORS request")
			wrappedGrpc.ServeHTTP(resp, req)
			return
		}

		if wrappedGrpc.IsGrpcWebRequest(req) {
			logger.Debug().Msg("is a gRPC-Web request")
			wrappedGrpc.ServeHTTP(resp, req)
			return
		}

		// Fall back to other servers.
		logger.Debug().Msg("s NOT a gRPC-Web or a gRPC-Web CORS request")
		http.DefaultServeMux.ServeHTTP(resp, req)
	})

	return &http.Server{
		Addr:           address,
		Handler:        grpcWebHandler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

func (s *Server) Register(service Service) {
	s.grpcServer.RegisterService(service.RPCServiceDesc(), service)
}

// Run starts the gRPC server using the configured listener
func (s *Server) Run() {
	if s.grpcWebServer != nil {
		go s.runGrpcWebServer()
	}

	s.logger.Info().
		Str("address", s.grpcAddress).
		Msg("Starting gRPC server...")

	err := s.grpcServer.Serve(s.grpcListener)
	if err != nil {
		s.logger.Fatal().Err(err).Msg("Error running gRPC server")
	}

	s.logger.Info().Msg("gRPC server has shut down")
}

func (s *Server) runGrpcWebServer() {
	s.logger.Info().
		Str("address", s.grpcWebAddress).
		Msg("Starting gRPC-Web server...")

	err := s.grpcWebServer.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("Error running gRPC-Web server")
	}

	s.logger.Info().Msg("gRPC-Web server has shut down")
}

func (s *Server) Close() error {
	if s.grpcWebServer != nil {
		s.logger.Info().Msg("Shutting down gRPC-Web server...")
		s.grpcWebServer.Close()
	}

	s.logger.Info().Msg("Shutting down gRPC server...")
	s.grpcServer.Stop()

	return nil
}
