package rpc

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

// Creates a test gRPC server. A gRPC-Web proxy is not created
func NewTestServer() (*Server, *grpc.ClientConn, error) {
	logger := log.Logger.With().Str("component", "rpc").Logger()

	logger.Info().Msg("Initializing test gRPC server...")

	listener := bufconn.Listen(10)

	server := &Server{
		config: Config{
			Port:    7101,
			WebPort: 0,
		},
		grpcServer:    grpc.NewServer(),
		grpcListener:  listener,
		grpcWebServer: nil,
		logger:        logger}

	return server, dialWithListener(listener), nil
}

func dialWithListener(listener *bufconn.Listener) *grpc.ClientConn {
	clientTarget := "localhost:11111" // probably nothing listens on this port, but that shouldn't matter because we'll dial a connection with the listener
	dialOptions := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithContextDialer(func(ctx context.Context, url string) (net.Conn, error) {
			return listener.Dial()
		})}

	conn, err := grpc.Dial(clientTarget, dialOptions...)
	if err != nil {
		panic(err)
	}

	return conn
}
