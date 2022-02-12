package command

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"

	pastebin_pb "github.com/BinacsLee/server/api/pastebin"
)

var (
	PastebinCmd = &cobra.Command{
		Use:   "pastebin",
		Short: "PasteBin Command:\t Just run `cli pastebin submit sth.(file)`",
		Args: func(cmd *cobra.Command, args []string) error {
			if !checkArgs(args, 2, 2) {
				return fmt.Errorf("error args length")
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			op, file := parsePastebinAgrs(args)
			switch op {
			case "submit":
				file, data := processReadFile(file)
				if len(file) != 0 {
					handleResp(node.Pastebin.PastebinSubmit(context.Background(), &pastebin_pb.PastebinSubmitReq{
						Text: string(data),
					}))
				} else {
					log.Printf(errorReadFile, file, data)
				}
			// TODO:
			// case "posts":
			default:
				log.Printf(errorOpInvalid)
			}
		},
	}
)

func parsePastebinAgrs(args []string) (op, file string) {
	return strings.ToLower(args[0]), args[1]
}
