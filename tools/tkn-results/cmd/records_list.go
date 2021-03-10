package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	pb "github.com/tektoncd/results/proto/v1alpha2/results_go_proto"
)

var (
	recordsListCmd = &cobra.Command{
		Use:   "list",
		Short: "List Records",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := client(ctx)
			if err != nil {
				return err
			}
			resp, err := client.ListRecords(ctx, &pb.ListRecordsRequest{
				Parent:    args[0],
				Filter:    filter,
				PageSize:  limit,
				PageToken: pageToken,
			})
			if err != nil {
				fmt.Printf("ListRecords: %v\n", err)
				return err
			}
			return printproto(os.Stdout, resp, format)
		},
		Args: cobra.ExactArgs(1),
	}
)

func init() {
	listFlags(recordsCmd.Flags())
	recordsCmd.AddCommand(recordsListCmd)
}
