package rpc

import (
	"google.golang.org/grpc"
)

type (
	RPCService struct {
		ServiceDescriptor grpc.ServiceDesc
	}

	Service interface {
		RPCServiceDesc() *grpc.ServiceDesc
	}
)

func (s RPCService) RPCServiceDesc() *grpc.ServiceDesc {
	return &s.ServiceDescriptor
}
