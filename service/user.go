package service

import (
	"context"

	"google.golang.org/grpc"

	user_pb "github.com/binacs/server/api/user"
)

// UserClientImpl Web tinyurl service implement
type UserClientImpl struct {
	Conn *grpc.ClientConn `inject-name:"Conn"`
	cli  user_pb.UserClient
}

// AfterInject do inject
func (impl *UserClientImpl) AfterInject() error {
	impl.cli = user_pb.NewUserClient(impl.Conn)
	return nil
}

// UserTest test
func (impl *UserClientImpl) UserTest(ctx context.Context, req *user_pb.UserTestReq) (*user_pb.UserTestResp, error) {
	return impl.cli.UserTest(ctx, req)
}

// UserRegister register
func (impl *UserClientImpl) UserRegister(ctx context.Context, req *user_pb.UserRegisterReq) (*user_pb.UserRegisterResp, error) {
	return impl.cli.UserRegister(ctx, req)
}

// UserAuth auth
func (impl *UserClientImpl) UserAuth(ctx context.Context, req *user_pb.UserAuthReq) (*user_pb.UserAuthResp, error) {
	return impl.cli.UserAuth(ctx, req)
}

// UserRefresh refresh
func (impl *UserClientImpl) UserRefresh(ctx context.Context, req *user_pb.UserRefreshReq) (*user_pb.UserRefreshResp, error) {
	return impl.cli.UserRefresh(ctx, req)
}

// UserInfo info
func (impl *UserClientImpl) UserInfo(ctx context.Context, req *user_pb.UserInfoReq) (*user_pb.UserInfoResp, error) {
	return impl.cli.UserInfo(ctx, req)
}
