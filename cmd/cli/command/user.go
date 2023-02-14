package command

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"

	user_pb "github.com/binacs/server/api/user"
)

var (
	UserCmd = &cobra.Command{
		Use:   "user",
		Short: "User Command:\t Just run `cli user test/register/auth/refresh/info`",
		Args: func(cmd *cobra.Command, args []string) error {
			if !checkArgs(args, 0, 0) {
				return fmt.Errorf("error args length")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			op, arg := parseUserAgrs(args)
			switch op {
			case "test":
				handleResp(node.User.UserTest(context.Background(), &user_pb.UserTestReq{}))
			case "register":
				if checkArgLength(arg, 2) {
					handleResp(node.User.UserRegister(context.Background(), &user_pb.UserRegisterReq{
						Id:  arg[0],
						Pwd: arg[1],
					}))
				}
			case "auth":
				if checkArgLength(arg, 2) {
					handleResp(node.User.UserAuth(context.Background(), &user_pb.UserAuthReq{
						Id:  arg[0],
						Pwd: arg[1],
					}))
				}
			case "refresh":
				if checkArgLength(arg, 1) {
					handleResp(node.User.UserRefresh(context.Background(), &user_pb.UserRefreshReq{
						RefreshToken: arg[0],
					}))
				}
			case "info":
				if checkArgLength(arg, 1) {
					handleResp(node.User.UserInfo(context.Background(), &user_pb.UserInfoReq{
						AccessToken: arg[0],
					}))
				}
			default:
				log.Printf(errorOpInvalid)
			}
		},
	}
)

func parseUserAgrs(args []string) (op string, arg []string) {
	return strings.ToLower(args[0]), args[1:]
}

func checkArgLength(arg []string, theshould int) bool {
	if len(arg) != theshould {
		log.Printf(errorArgsLengthInvalid, theshould, len(arg))
		return false
	}
	return true
}
