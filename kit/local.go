package kit

import (
	"github.com/caspr-io/mu-kit/rpc"
	"github.com/caspr-io/mu-kit/streaming"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func NewLocalTestKitServer(name string, f func(*rpc.Server, *streaming.River) rpc.Service) (*MuKitServer, *grpc.ClientConn) {
	initLogger(name)
	log.Logger.Info().Msg("Starting local µ-Kit server...")

	river, err := streaming.NewTestRiver()
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("Failed to initialize µ-Kit Streaming system")
	}

	listener := bufconn.Listen(10)

	rpcServer, err := rpc.CreateServer(listener)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("Failed to initialize µ-Kit gRPC server")
	}

	rpcService := f(rpcServer, river)
	rpcServer.Register(rpcService)

	server := CreateKit(&MuKitConfig{}, rpcServer, river)

	return server, rpc.DialConnection(listener)
}
