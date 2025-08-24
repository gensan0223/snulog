package test

import (
	"context"
	"testing"
	"time"

	pb "github.com/gensan0223/snulog/proto"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TestGRPCConnection tests if we can connect to the gRPC server
func TestGRPCConnection(t *testing.T) {
	// Skip this test if no gRPC server is running
	t.Skip("Skipping integration test - requires running gRPC server")

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			t.Logf("Failed to close connection: %v", err)
		}
	}()

	client := pb.NewLogServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Test AddLogs
	entry := &pb.LogEntry{
		UserName:  "integration_test",
		Status:    "testing",
		Feeling:   "🧪",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	resp, err := client.AddLogs(ctx, entry)
	if err != nil {
		t.Fatalf("AddLogs failed: %v", err)
	}

	if resp.Message == "" {
		t.Error("Expected non-empty response message")
	}

	// Test FetchLogs
	fetchResp, err := client.FetchLogs(ctx, &pb.FetchRequest{
		TeamId: "default",
	})
	if err != nil {
		t.Fatalf("FetchLogs failed: %v", err)
	}

	if len(fetchResp.Logs) == 0 {
		t.Error("Expected at least one log entry")
	}
}

// TestWebHandlerWithMockServer tests the web handler with a mock gRPC server
func TestWebHandlerWithMockServer(t *testing.T) {
	// Skip this test if no database is available
	t.Skip("Skipping test - requires database connection")

	// This would require setting up a mock gRPC server and database
	// For now, we'll skip this test
}

// TestWebEndpoints tests the web endpoints
func TestWebEndpoints(t *testing.T) {
	// Skip this test if no database is available
	t.Skip("Skipping integration test - requires database connection")

	// In a real test environment, you would set up a test database
	// For example:
	// db, err := sql.Open("postgres", "postgres://test:test@localhost:5432/test_db?sslmode=disable")
	// if err != nil {
	//     t.Fatalf("Failed to connect to test database: %v", err)
	// }
	// defer db.Close()
	//
	// webHandler := handler.NewWebHandler("localhost:50051", db)
	// ... rest of the test
}
