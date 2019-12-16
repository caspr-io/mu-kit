package kit

import (
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/caspr-io/mu-kit/rpc"
	"github.com/caspr-io/mu-kit/streaming"
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
