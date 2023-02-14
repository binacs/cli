package service

import (
	"context"

	cos_pb "github.com/binacs/server/api/cos"
	crypto_pb "github.com/binacs/server/api/crypto"
	pastebin_pb "github.com/binacs/server/api/pastebin"
	tinyurl_pb "github.com/binacs/server/api/tinyurl"
	user_pb "github.com/binacs/server/api/user"
)

// CosClient tinyurl client
type CosClient interface {
	CosBucketURL(ctx context.Context, req *cos_pb.CosBucketURLReq) (*cos_pb.CosBucketURLResp, error)
	CosPut(ctx context.Context, req *cos_pb.CosPutReq) (*cos_pb.CosPutResp, error)
	CosGet(ctx context.Context, req *cos_pb.CosGetReq) (*cos_pb.CosGetResp, error)
}

// CryptoClient crypto client
type CryptoClient interface {
	CryptoEncrypt(ctx context.Context, req *crypto_pb.CryptoEncryptReq) (*crypto_pb.CryptoEncryptResp, error)
	CryptoDecrypt(ctx context.Context, req *crypto_pb.CryptoDecryptReq) (*crypto_pb.CryptoDecryptResp, error)
}

// PastebinClient pastebin client
type PastebinClient interface {
	PastebinSubmit(ctx context.Context, req *pastebin_pb.PastebinSubmitReq) (*pastebin_pb.PastebinSubmitResp, error)
}

// TinyURLClient tinyurl client
type TinyURLClient interface {
	TinyURLEncode(ctx context.Context, req *tinyurl_pb.TinyURLEncodeReq) (*tinyurl_pb.TinyURLEncodeResp, error)
	TinyURLDecode(ctx context.Context, req *tinyurl_pb.TinyURLDecodeReq) (*tinyurl_pb.TinyURLDecodeResp, error)
}

// UserClient user client
type UserClient interface {
	UserTest(ctx context.Context, req *user_pb.UserTestReq) (*user_pb.UserTestResp, error)
	UserRegister(ctx context.Context, req *user_pb.UserRegisterReq) (*user_pb.UserRegisterResp, error)
	UserAuth(ctx context.Context, req *user_pb.UserAuthReq) (*user_pb.UserAuthResp, error)
	UserRefresh(ctx context.Context, req *user_pb.UserRefreshReq) (*user_pb.UserRefreshResp, error)
	UserInfo(ctx context.Context, req *user_pb.UserInfoReq) (*user_pb.UserInfoResp, error)
}
