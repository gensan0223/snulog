package cmd

import (
	"context"
	"fmt"
	"time"

	pb "github.com/gensan0223/snulog/proto"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new progress and emotion log",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 3 {
			fmt.Println("引数が足りません")
		}

		type LogEntry struct {
			UserName  string
			Status    string
			Feeling   string
			Timestamp string
		}

		entry := &pb.LogEntry{
			UserName:  args[0],
			Status:    args[1],
			Feeling:   args[2],
			Timestamp: time.Now().Format(time.RFC3339),
		}

		conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Println("⛔gRPC接続失敗: ", err)
			return
		}
		defer conn.Close()

		conn.Connect()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		client := pb.NewLogServiceClient(conn)
		res, err := client.AddLogs(ctx, entry)
		if err != nil {
			fmt.Println("⛔ログ追加失敗: ", err)
			return
		}

		fmt.Printf("✅ログ追加 \nuser: %s\nstatus: %s\nfeeling: %s\ntimestamp: %s\n", entry.UserName, entry.Status, entry.Feeling, entry.Timestamp)
		fmt.Printf("✅サーバ応答: %s\n", res.Message)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
