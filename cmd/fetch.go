package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/gensan0223/snulog/internal/util"
	pb "github.com/gensan0223/snulog/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "チームメンバーの進捗と感情ログを取得する",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Println("⛔gRPC接続失敗: ", err)
			return
		}
		defer util.CloseWithLog(conn)

		client := pb.NewLogServiceClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		resp, err := client.FetchLogs(ctx, &pb.FetchRequest{
			TeamId: "default",
		})
		if err != nil {
			fmt.Println("⛔データ取得失敗: ", err)
			return
		}
		for _, log := range resp.Logs {
			fmt.Printf("👤 %s\t📝 %s\t😀 %s\t🕒 %s\n", log.UserName, log.Status, log.Feeling, log.Timestamp)
		}
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
