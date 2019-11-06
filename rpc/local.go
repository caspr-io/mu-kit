package rpc

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func DialConnection(listener *bufconn.Listener) *grpc.ClientConn {
	dialOptions := []grpc.DialOption{grpc.WithInsecure(), grpc.WithContextDialer(getBufconnDialer(listener))}
	clientTarget := "localhost:11111" // probably nothing listens on this port, but that shouldn't matter because we'll use the bufconn dialer

	conn, err := grpc.Dial(clientTarget, dialOptions...)
	if err != nil {
		panic(err)
	}

	return conn

}

func getBufconnDialer(listener *bufconn.Listener) func(context.Context, string) (net.Conn, error) {
	return func(ctx context.Context, url string) (net.Conn, error) {
		return listener.Dial()
	}
}
