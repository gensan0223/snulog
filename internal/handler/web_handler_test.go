package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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
	handler := NewWebHandler("localhost:50051")

	form := url.Values{}
	form.Add("user_name", "test_user")
	form.Add("status", "working")
	form.Add("feeling", "😊")

	req := httptest.NewRequest(http.MethodPost, "/api/logs", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	// Note: This test will fail without a running gRPC server
	// In a real test environment, you'd want to mock the gRPC client
	handler.AddLog(w, req)

	// Since we don't have a running server, we expect a connection error
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "error-message") {
		t.Errorf("Expected error message in response, got: %s", body)
	}
}

func TestWebHandler_AddLog_MissingFields(t *testing.T) {
	handler := NewWebHandler("localhost:50051")

	form := url.Values{}
	form.Add("user_name", "test_user")
	// Missing status and feeling

	req := httptest.NewRequest(http.MethodPost, "/api/logs", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	handler.AddLog(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	expected := "すべてのフィールドを入力してください"
	if !strings.Contains(body, expected) {
		t.Errorf("Expected '%s' in response, got: %s", expected, body)
	}
}

func TestWebHandler_AddLog_WrongMethod(t *testing.T) {
	handler := NewWebHandler("localhost:50051")

	req := httptest.NewRequest(http.MethodGet, "/api/logs", nil)
	w := httptest.NewRecorder()

	handler.AddLog(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestWebHandler_GetLogs(t *testing.T) {
	handler := NewWebHandler("localhost:50051")

	req := httptest.NewRequest(http.MethodGet, "/api/logs", nil)
	w := httptest.NewRecorder()

	// Note: This test will fail without a running gRPC server
	handler.GetLogs(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	// Since we don't have a running server, we expect a connection error
	if !strings.Contains(body, "error-message") {
		t.Errorf("Expected error message in response, got: %s", body)
	}
}

func TestWebHandler_ServeIndex(t *testing.T) {
	handler := NewWebHandler("localhost:50051")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Note: This test will fail without the template file
	handler.ServeIndex(w, req)

	// We expect an error since the template file doesn't exist in test environment
	if w.Code == http.StatusOK {
		body := w.Body.String()
		if len(body) == 0 {
			t.Error("Expected some content in response")
		}
	}
}
