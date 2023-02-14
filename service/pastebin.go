package service

import (
	"context"

	"google.golang.org/grpc"

	pastebin_pb "github.com/binacs/server/api/pastebin"
)

// PastebinClientImpl Web crypto service implement
type PastebinClientImpl struct {
	Conn *grpc.ClientConn `inject-name:"Conn"`
	cli  pastebin_pb.PastebinClient
}

// AfterInject do inject
func (impl *PastebinClientImpl) AfterInject() error {
	impl.cli = pastebin_pb.NewPastebinClient(impl.Conn)
	return nil
}

// PastebinSubmit submit
func (impl *PastebinClientImpl) PastebinSubmit(ctx context.Context, req *pastebin_pb.PastebinSubmitReq) (*pastebin_pb.PastebinSubmitResp, error) {
	return impl.cli.PastebinSubmit(ctx, req)
}
