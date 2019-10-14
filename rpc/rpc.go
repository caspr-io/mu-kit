package rpc

import (
	"google.golang.org/grpc"
)

type (
	RpcService struct {
		ServiceDescriptor grpc.ServiceDesc
	}
)

func (s RpcService) RPCServiceDesc() *grpc.ServiceDesc {
	return &s.ServiceDescriptor
}
