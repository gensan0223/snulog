package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gensan0223/snulog/internal/handler"
	pb "github.com/gensan0223/snulog/proto"
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
	defer conn.Close()

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
	// This would require setting up a mock gRPC server
	// For now, we'll test the handler logic without the server

	handler := handler.NewWebHandler("localhost:50051")

	// Test form validation
	form := url.Values{}
	form.Add("user_name", "test")
	form.Add("status", "working")
	form.Add("feeling", "😊")

	req := httptest.NewRequest(http.MethodPost, "/api/logs", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	handler.AddLog(w, req)

	// We expect this to fail with connection error since no server is running
	body := w.Body.String()
	if !strings.Contains(body, "error-message") {
		t.Errorf("Expected error message due to no server, got: %s", body)
	}
}

// TestWebEndpoints tests the web endpoints
func TestWebEndpoints(t *testing.T) {
	webHandler := handler.NewWebHandler("localhost:50051")

	// Create a test server
	mux := http.NewServeMux()
	mux.HandleFunc("/", webHandler.ServeIndex)
	mux.HandleFunc("/api/logs", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			webHandler.AddLog(w, r)
		} else if r.Method == http.MethodGet {
			webHandler.GetLogs(w, r)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	// Test GET /api/logs
	resp, err := http.Get(server.URL + "/api/logs")
	if err != nil {
		t.Fatalf("Failed to GET /api/logs: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Test POST /api/logs
	form := url.Values{}
	form.Add("user_name", "test")
	form.Add("status", "working")
	form.Add("feeling", "😊")

	resp, err = http.PostForm(server.URL+"/api/logs", form)
	if err != nil {
		t.Fatalf("Failed to POST /api/logs: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}
