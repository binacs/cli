package command

import (
	"context"
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"github.com/binacs/server/types"

	cos_pb "github.com/binacs/server/api/cos"
	crypto_pb "github.com/binacs/server/api/crypto"
	pastebin_pb "github.com/binacs/server/api/pastebin"
	tinyurl_pb "github.com/binacs/server/api/tinyurl"
	user_pb "github.com/binacs/server/api/user"

	"github.com/binacs/cli/service"
	"github.com/binacs/cli/util"
)

var instance, domain, port string

var (
	StartCmd = &cobra.Command{
		Use:   "start",
		Short: "Start Command",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			setInstance()
			return dailAndServe()
		},
	}
)

func init() {
	startCmdFlags(StartCmd)
}

func startCmdFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&instance, "instance", "", "Local instance name")
	cmd.PersistentFlags().StringVar(&domain, "domain", "api.binacs.space", "API domain such as api.binacs.space")
	cmd.PersistentFlags().StringVar(&port, "port", ":30000", "API port such as :30000")
}

func setInstance() {
	if len(instance) == 0 {
		hostname, err := os.Hostname()
		if err != nil {
			log.Printf("os.Hostname get err: %+v", err)
			instance = "defaultInstanceName"
		} else {
			instance = hostname
		}
	}
	log.Printf("instance = %s domain = %s port = %s", instance, domain, port)
}

func dailAndServe() error {
	// Dail API server. No long-lived credential is configured here: the
	// server never actually validates the bearer token for non-Cos RPCs
	// (see server/gateway/grpc.go), so a fixed placeholder keeps them
	// working; Cos RPCs instead carry a per-call pass key relayed from an
	// interactive prompt in `cli` (see util.RelayAuth / util.AttachAuth) —
	// nothing sensitive is ever persisted to disk by clid.
	conn, err := grpc.Dial(domain+port,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(util.GetCertPool(), domain)),
		grpc.WithChainUnaryInterceptor(defaultAuthInterceptor),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(types.GrpcMsgSize)),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(types.GrpcMsgSize)),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Listen to unix socket
	_ = os.Remove(util.GetSockPath())
	sockAddr, err := net.ResolveUnixAddr("unix", util.GetSockPath())
	if err != nil {
		return err
	}
	lis, err := net.ListenUnix("unix", sockAddr)
	if err != nil {
		return err
	}

	// Serve
	s := grpc.NewServer(
		grpc.MaxRecvMsgSize(types.GrpcMsgSize),
		grpc.MaxSendMsgSize(types.GrpcMsgSize),
	)
	node := service.InitService(conn)

	cos_pb.RegisterCosServer(s, node.Cos.(cos_pb.CosServer))
	crypto_pb.RegisterCryptoServer(s, node.Crypto.(crypto_pb.CryptoServer))
	pastebin_pb.RegisterPastebinServer(s, node.Pastebin.(pastebin_pb.PastebinServer))
	tinyurl_pb.RegisterTinyURLServer(s, node.TinyURL.(tinyurl_pb.TinyURLServer))
	user_pb.RegisterUserServer(s, node.User.(user_pb.UserServer))

	return s.Serve(lis)
}

// defaultAuthInterceptor stamps a placeholder authorization header onto
// outbound calls that don't already carry one (e.g. Cos calls that had a
// real pass key relayed onto their outgoing context via util.RelayAuth).
// This exists only so the server's blanket "header must be present" check
// doesn't reject unrelated services; it is not a credential.
func defaultAuthInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	if md, ok := metadata.FromOutgoingContext(ctx); !ok || len(md.Get(util.HeaderAuthorize)) == 0 {
		ctx = metadata.AppendToOutgoingContext(ctx, util.HeaderAuthorize, util.TokenPrefix+"unused")
	}
	return invoker(ctx, method, req, reply, cc, opts...)
}
