package service

import (
	"context"

	"google.golang.org/grpc"

	crypto_pb "github.com/binacs/server/api/crypto"
)

// CryptoClientImpl Web crypto client implement
type CryptoClientImpl struct {
	Conn *grpc.ClientConn `inject-name:"Conn"`
	cli  crypto_pb.CryptoClient
}

// AfterInject do inject
func (impl *CryptoClientImpl) AfterInject() error {
	impl.cli = crypto_pb.NewCryptoClient(impl.Conn)
	return nil
}

// CryptoEncrypt encrypt
func (impl *CryptoClientImpl) CryptoEncrypt(ctx context.Context, req *crypto_pb.CryptoEncryptReq) (*crypto_pb.CryptoEncryptResp, error) {
	return impl.cli.CryptoEncrypt(ctx, req)
}

// CryptoDecrypt decrypt
func (impl *CryptoClientImpl) CryptoDecrypt(ctx context.Context, req *crypto_pb.CryptoDecryptReq) (*crypto_pb.CryptoDecryptResp, error) {
	return impl.cli.CryptoDecrypt(ctx, req)
}
