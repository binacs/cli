package command

import (
	"log"
	"net"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

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
	cmd.PersistentFlags().StringVar(&domain, "domain", "api.binacs.cn", "API domain such as api.binacs.cn")
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
	// Dail API server
	conn, err := grpc.Dial(domain+port,
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(util.GetCertPool(), domain)),
		grpc.WithPerRPCCredentials(util.GetToken(instance)),
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
