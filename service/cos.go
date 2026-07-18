package service

import (
	"context"

	"google.golang.org/grpc"

	cos_pb "github.com/binacs/server/api/cos"

	"github.com/binacs/cli/util"
)

// CosClientImpl Web crypto client implement
type CosClientImpl struct {
	Conn *grpc.ClientConn `inject-name:"Conn"`
	cli  cos_pb.CosClient
}

// AfterInject do inject
func (impl *CosClientImpl) AfterInject() error {
	impl.cli = cos_pb.NewCosClient(impl.Conn)
	return nil
}

// CosBucketURL return the bucket url
func (impl *CosClientImpl) CosBucketURL(ctx context.Context, req *cos_pb.CosBucketURLReq) (*cos_pb.CosBucketURLResp, error) {
	return impl.cli.CosBucketURL(ctx, req)
}

// CosPut put. The pass key cli attached on the local socket call
// (util.AttachAuth) is relayed as-is onto the upstream call to the
// server; clid never sees or stores it beyond this one request.
func (impl *CosClientImpl) CosPut(ctx context.Context, req *cos_pb.CosPutReq) (*cos_pb.CosPutResp, error) {
	return impl.cli.CosPut(util.RelayAuth(ctx), req)
}

// CosPut get
func (impl *CosClientImpl) CosGet(ctx context.Context, req *cos_pb.CosGetReq) (*cos_pb.CosGetResp, error) {
	return impl.cli.CosGet(util.RelayAuth(ctx), req)
}
