package grpc

import (
	"context"
	"net"

	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var (
	lis   *bufconn.Listener
	err   error
	zero  uint32 = 0
	one   uint32 = 1
	two   uint32 = 2
	three uint32 = 3
)

func bufDialer(ctx context.Context, address string) (net.Conn, error) {
	return lis.Dial()
}
