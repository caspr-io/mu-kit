package test

import (
	"context"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

func GetClientConn(startLocalServer func(listener *bufconn.Listener)) (listener *grpc.ClientConn) {
	dialOptions := []grpc.DialOption{grpc.WithInsecure()}

	clientTarget := os.Getenv("MUKIT_INTEGRATION_TEST_TARGET")
	if clientTarget == "" {
		clientTarget = "localhost:11111" // probably nothing listens on this port, but that shouldn't matter because we'll use the bufconn dialer
		listener := bufconn.Listen(10)
		go startLocalServer(listener)
		dialOptions = []grpc.DialOption{grpc.WithInsecure(), grpc.WithContextDialer(getBufconnDialer(listener))}
	}

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
