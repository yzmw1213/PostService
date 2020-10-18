package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var (
	lis  *bufconn.Listener
	err  error
	zero int32 = 0
	one  int32 = 1
)

func bufDialer(ctx context.Context, address string) (net.Conn, error) {
	return lis.Dial()
}
