package kit

import (
	"github.com/caspr-io/mu-kit/rpc"
	"github.com/caspr-io/mu-kit/streaming"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func NewTestSystem() (*MuKitServer, *grpc.ClientConn) {
	log.Logger.Info().Msg("Starting local µ-Kit server...")

	river, err := streaming.NewTestRiver()
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("Cannot initialize local µ-Kit Streaming system")
	}

	rpcServer, clientConn, err := rpc.NewTestServer()
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("Cannot initialize local µ-Kit gRPC server")
	}

	server := createSystem(&MuKitConfig{}, rpcServer, river)

	return server, clientConn
}

// Deprecated: use kit.InitLogger() and kit.NewTestSystem() instead
func NewLocalTestKitServer(name string, f func(*rpc.Server, *streaming.River) rpc.Service) (*MuKitServer, *grpc.ClientConn) {
	InitLogger(name)

	server, conn := NewTestSystem()

	rpcService := f(server.RPCServer(), server.River())
	server.RPCServer().Register(rpcService)

	return server, conn
}
