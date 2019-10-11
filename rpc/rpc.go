package rpc

import (
	"net/http"

	"github.com/NYTimes/gizmo/server/kit"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"google.golang.org/grpc"
)

type (
	RpcService struct {
		ServiceDescription grpc.ServiceDesc
	}
)

func (s RpcService) HTTPEndpoints() map[string]map[string]kit.HTTPEndpoint {
	return nil
}

func (s RpcService) HTTPMiddleware(h http.Handler) http.Handler {
	return nil
}

func (s RpcService) HTTPOptions() []httptransport.ServerOption {
	return nil
}

func (s RpcService) HTTPRouterOptions() []kit.RouterOption {
	return nil
}

func (s RpcService) Middleware(e endpoint.Endpoint) endpoint.Endpoint {
	return e
}

func (s RpcService) RPCMiddleware() grpc.UnaryServerInterceptor {
	return nil
}

func (s RpcService) RPCOptions() []grpc.ServerOption {
	return nil
}

func (s RpcService) RPCServiceDesc() *grpc.ServiceDesc {
	// snagged from the pb.go file
	return &s.ServiceDescription
}
