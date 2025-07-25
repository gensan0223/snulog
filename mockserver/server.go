package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/gensan0223/snulog/proto"

	"google.golang.org/grpc"
)

type logServer struct {
	pb.UnimplementedLogServiceServer
	logs []*pb.LogEntry
}

func (s *logServer) AddLogs(ctx context.Context, entry *pb.LogEntry) (*pb.AddResponse, error) {
	s.logs = append(s.logs, entry)
	log.Printf("Received AddLog entry: %+v", entry)

	return &pb.AddResponse{Message: "added successfully"}, nil
}

func (s *logServer) FetchLogs(ctx context.Context, req *pb.FetchRequest) (*pb.FetchResponse, error) {
	log.Printf("Received FetchLogs request")

	return &pb.FetchResponse{
		Logs: s.logs,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := &logServer{}

	grpcServer := grpc.NewServer()
	pb.RegisterLogServiceServer(grpcServer, srv)
	fmt.Printf("✅ Mock gRPC server listening on %s", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
