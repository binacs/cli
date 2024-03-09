package command

import (
	"net"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"

	"github.com/binacs/server/types"

	"github.com/binacs/cli/service"
	"github.com/binacs/cli/util"
)

var node *service.NodeServiceImpl

func unixConnect(addr string, t time.Duration) (net.Conn, error) {
	unix_addr, err := net.ResolveUnixAddr("unix", util.GetSockPath())
	conn, err := net.DialUnix("unix", nil, unix_addr)
	return conn, err
}

var (
	RootCmd = &cobra.Command{
		Use:   "root",
		Short: "Terminal client for https://binacs.space\nMore at https://github.com/binacs/cli",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			conn, err := grpc.Dial(util.GetSockPath(),
				grpc.WithBlock(),
				grpc.WithInsecure(),
				grpc.WithDialer(unixConnect),
				grpc.FailOnNonTempDialError(true),
				grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(types.GrpcMsgSize)),
				grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(types.GrpcMsgSize)),
			)
			if err != nil {
				return err
			}
			node = service.InitService(conn)
			return nil
		},
	}
)
