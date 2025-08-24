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

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Debug gRPC connection",
	Run: func(cmd *cobra.Command, args []string) {
		grpcAddr := "localhost:50051"

		fmt.Printf("🔍 Attempting to connect to gRPC server at %s\n", grpcAddr)

		conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("❌ Failed to create gRPC client: %v\n", err)
			return
		}
		defer func() {
			if err := conn.Close(); err != nil {
				fmt.Printf("⚠️ Failed to close connection: %v\n", err)
			}
		}()

		fmt.Println("✅ gRPC client created successfully")

		client := pb.NewLogServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		// Test AddLogs
		fmt.Println("🧪 Testing AddLogs...")
		entry := &pb.LogEntry{
			UserName:  "debug_user",
			Status:    "testing",
			Feeling:   "🔧",
			Timestamp: time.Now().Format(time.RFC3339),
		}

		resp, err := client.AddLogs(ctx, entry)
		if err != nil {
			fmt.Printf("❌ AddLogs failed: %v\n", err)
		} else {
			fmt.Printf("✅ AddLogs succeeded: %s\n", resp.Message)
		}

		// Test FetchLogs
		fmt.Println("🧪 Testing FetchLogs...")
		fetchResp, err := client.FetchLogs(ctx, &pb.FetchRequest{
			TeamId: "default",
		})
		if err != nil {
			fmt.Printf("❌ FetchLogs failed: %v\n", err)
		} else {
			fmt.Printf("✅ FetchLogs succeeded, got %d logs\n", len(fetchResp.Logs))
			for i, log := range fetchResp.Logs {
				fmt.Printf("  %d: %s - %s - %s (%s)\n", i+1, log.UserName, log.Status, log.Feeling, log.Timestamp)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)
}
