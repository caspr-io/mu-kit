package kit

import (
	"github.com/caspr-io/mu-kit/river"
	"github.com/caspr-io/mu-kit/rpc"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func NewLocalTestKitServer(f func(*rpc.SubSystem, *river.SubSystem) rpc.RPCService) (*MuKitServer, *grpc.ClientConn) {
	initLogger()
	log.Logger.Info().Msg("Starting local µ-Kit server..")

	riverSystem, err := river.NewTestRiver()
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("Failed to start River subsystem")
	}

	listener := bufconn.Listen(10)

	rpcSystem, err := rpc.NewWithListener(listener)
	if err != nil {
		log.Logger.Fatal().Err(err).Msg("Failed to start RPC subsystem")
	}

	rpcService := f(rpcSystem, riverSystem)
	rpcSystem.Register(rpcService)

	server := NewWithSubSystems(rpcSystem, riverSystem)

	return server, rpc.DialConnection(listener)
}
