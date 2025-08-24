package handler

import (
	"context"
	"testing"

	pb "github.com/gensan0223/snulog/proto"
	"google.golang.org/grpc"
)

// Mock gRPC client for testing
type mockLogServiceClient struct {
	addLogsFunc   func(ctx context.Context, in *pb.LogEntry, opts ...grpc.CallOption) (*pb.AddResponse, error)
	fetchLogsFunc func(ctx context.Context, in *pb.FetchRequest, opts ...grpc.CallOption) (*pb.FetchResponse, error)
}

func (m *mockLogServiceClient) AddLogs(ctx context.Context, in *pb.LogEntry, opts ...grpc.CallOption) (*pb.AddResponse, error) {
	if m.addLogsFunc != nil {
		return m.addLogsFunc(ctx, in, opts...)
	}
	return &pb.AddResponse{Message: "success"}, nil
}

func (m *mockLogServiceClient) FetchLogs(ctx context.Context, in *pb.FetchRequest, opts ...grpc.CallOption) (*pb.FetchResponse, error) {
	if m.fetchLogsFunc != nil {
		return m.fetchLogsFunc(ctx, in, opts...)
	}
	return &pb.FetchResponse{
		Logs: []*pb.LogEntry{
			{
				UserName:  "test_user",
				Status:    "working",
				Feeling:   "😊",
				Timestamp: "2024-01-01T00:00:00Z",
			},
		},
	}, nil
}

func TestWebHandler_AddLog_Success(t *testing.T) {
	// Skip this test since it requires a database connection
	t.Skip("Skipping test - requires database connection and running gRPC server")
}

func TestWebHandler_AddLog_MissingFields(t *testing.T) {
	// Skip this test since it requires a database connection
	t.Skip("Skipping test - requires database connection")
}

func TestWebHandler_AddLog_WrongMethod(t *testing.T) {
	// Skip this test since it requires a database connection
	t.Skip("Skipping test - requires database connection")
}

func TestWebHandler_GetLogs(t *testing.T) {
	// Skip this test since it requires a database connection
	t.Skip("Skipping test - requires database connection and running gRPC server")
}

func TestWebHandler_ServeIndex(t *testing.T) {
	// Skip this test since it requires a database connection
	t.Skip("Skipping test - requires database connection and template files")
}
