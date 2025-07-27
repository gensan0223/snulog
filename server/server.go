package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/gensan0223/snulog/internal/repository"
	"github.com/gensan0223/snulog/internal/usecase"
	"github.com/gensan0223/snulog/internal/util"
	pb "github.com/gensan0223/snulog/proto"

	_ "github.com/lib/pq"

	"google.golang.org/grpc"
)

type logServer struct {
	pb.UnimplementedLogServiceServer
	usecase usecase.LogUsecase
}

func (s *logServer) AddLogs(ctx context.Context, entry *pb.LogEntry) (*pb.AddResponse, error) {
	return s.usecase.AddLogs(ctx, entry)
}

func (s *logServer) FetchLogs(ctx context.Context, req *pb.FetchRequest) (*pb.FetchResponse, error) {
	return s.usecase.FetchLogs(ctx)
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer util.CloseWithLog(db)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// repo := repository.NewInMemoryLogRepository()
	repo := repository.NewPostgresLogRepository(db)
	uc := usecase.NewLogUsecase(repo)
	srv := &logServer{
		usecase: uc,
	}

	grpcServer := grpc.NewServer()
	pb.RegisterLogServiceServer(grpcServer, srv)
	fmt.Printf("✅ Mock gRPC server listening on %s", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
