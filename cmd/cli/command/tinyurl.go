package command

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"

	tinyurl_pb "github.com/binacs/server/api/tinyurl"
)

var (
	TinyurlCmd = &cobra.Command{
		Use:   "tinyurl",
		Short: "TinyURL Command:\t Just run `cli tinyurl encode/decode sth.`",
		Args: func(cmd *cobra.Command, args []string) error {
			if !checkArgs(args, 2, 2) {
				return fmt.Errorf("error args length")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			op, url := parseTinyurlAgrs(args)
			switch op {
			case "encode":
				handleResp(node.TinyURL.TinyURLEncode(context.Background(), &tinyurl_pb.TinyURLEncodeReq{
					Url: url,
				}))
			case "decode":
				handleResp(node.TinyURL.TinyURLDecode(context.Background(), &tinyurl_pb.TinyURLDecodeReq{
					Turl: url,
				}))
			default:
				log.Printf(errorOpInvalid)
			}
		},
	}
)

func parseTinyurlAgrs(args []string) (op, url string) {
	return strings.ToLower(args[0]), args[1]
}
