package service

import (
	"context"

	"google.golang.org/grpc"

	tinyurl_pb "github.com/binacs/server/api/tinyurl"
)

// TinyURLClientImpl client implement
type TinyURLClientImpl struct {
	Conn *grpc.ClientConn `inject-name:"Conn"`
	cli  tinyurl_pb.TinyURLClient
}

// AfterInject do inject
func (impl *TinyURLClientImpl) AfterInject() error {
	impl.cli = tinyurl_pb.NewTinyURLClient(impl.Conn)
	return nil
}

// TinyURLEncode encode
func (impl *TinyURLClientImpl) TinyURLEncode(ctx context.Context, req *tinyurl_pb.TinyURLEncodeReq) (*tinyurl_pb.TinyURLEncodeResp, error) {
	return impl.cli.TinyURLEncode(ctx, req)
}

// TinyURLDecode decode
func (impl *TinyURLClientImpl) TinyURLDecode(ctx context.Context, req *tinyurl_pb.TinyURLDecodeReq) (*tinyurl_pb.TinyURLDecodeResp, error) {
	return impl.cli.TinyURLDecode(ctx, req)
}
