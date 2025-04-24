package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/gensan0223/snulog/proto"

	"google.golang.org/grpc"
)

type mockLogServer struct {
	pb.UnimplementedLogServiceServer
}

func (s *mockLogServer) FetchLogs(ctx context.Context, req *pb.FetchRequest) (*pb.FetchResponse, error) {
	log.Printf("Received FetchLogs request for team: %s", req.TeamId)

	return &pb.FetchResponse{
		Logs: []*pb.LogEntry{
			{
				UserName:  "Alice",
				Status:    "コードレビュー中",
				Feeling:   "😀",
				Timestamp: "2025-08-01T10:00:00Z",
			},
			{
				UserName:  "Jane Doe",
				Status:    "Working",
				Feeling:   "🤮",
				Timestamp: "2023-08-01T10:00:00Z",
			},
		},
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterLogServiceServer(grpcServer, &mockLogServer{})
	fmt.Printf("✅ Mock gRPC server listening on %s", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
