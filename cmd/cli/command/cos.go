package command

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"

	cos_pb "github.com/binacs/server/api/cos"

	"github.com/binacs/cli/util"
)

var (
	CosCmd = &cobra.Command{
		Use:   "cos",
		Short: "Cos Command:\t Just run `cli cos put/get sth.(file)`",
		Args: func(cmd *cobra.Command, args []string) error {
			if !checkArgs(args, 2, 2) {
				return fmt.Errorf("error args length")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			op, arg := parseCosAgrs(args)
			switch op {
			case "put":
				file, data := processReadFile(arg)
				if len(file) == 0 {
					log.Printf(errorReadFile, file, data)
					return
				}
				ctx, err := cosAuthContext()
				if err != nil {
					log.Printf("Error: read pass key: %+v\n", err)
					return
				}
				handleResp(node.Cos.CosPut(ctx, &cos_pb.CosPutReq{
					FileName:  file,
					FileBytes: data,
				}))
			case "get":
				log.Printf("Error: Not support `get`.\n")
				// ctx, err := cosAuthContext()
				// handleResp(node.Cos.CosGet(ctx, &cos_pb.CosGetReq{
				// 	CosURI: arg,
				// }))
			default:
				log.Printf(errorOpInvalid)
			}
		},
	}
)

func parseCosAgrs(args []string) (op, arg string) {
	return strings.ToLower(args[0]), args[1]
}

// cosAuthContext prompts for the COS pass key on the controlling terminal
// (never persisted to disk) and attaches it to a fresh context for a
// single Cos RPC.
func cosAuthContext() (context.Context, error) {
	key, err := util.PromptSecret("COS pass key: ")
	if err != nil {
		return nil, err
	}
	return util.AttachAuth(context.Background(), key), nil
}
