package cmd

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/gensan0223/snulog/internal/util"
	"github.com/gensan0223/snulog/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

type mockLogServer struct {
	proto.UnimplementedLogServiceServer
}

func (s *mockLogServer) FetchLogs(ctx context.Context, req *proto.FetchRequest) (*proto.FetchResponse, error) {
	if req.TeamId == "default" {
		return &proto.FetchResponse{
			Logs: []*proto.LogEntry{
				{
					UserName:  "テストユーザー",
					Status:    "テスト中",
					Feeling:   "😊",
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	return nil, status.Error(codes.NotFound, "チームが見つかりません")
}

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	proto.RegisterLogServiceServer(s, &mockLogServer{})
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

func TestFetchLogs(t *testing.T) {
	tests := []struct {
		name        string
		teamID      string
		expectedErr bool
	}{
		{
			name:        "正常系: デフォルトチームのログ取得",
			teamID:      "default",
			expectedErr: false,
		},
		{
			name:        "異常系: 存在しないチームID",
			teamID:      "non-existent",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			dialer := func(context.Context, string) (net.Conn, error) {
				return lis.Dial()
			}
			// nolint:staticcheck // bufconn テスト用に grpc.DialContext を使用
			conn, err := grpc.DialContext(ctx, "bufnet",
				grpc.WithContextDialer(dialer),
				grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				t.Fatalf("gRPC接続に失敗しました: %v", err)
			}
			defer util.CloseWithLog(conn)

			client := proto.NewLogServiceClient(conn)

			resp, err := client.FetchLogs(ctx, &proto.FetchRequest{
				TeamId: tt.teamID,
			})

			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				if resp != nil {
					for _, log := range resp.Logs {
						assert.NotEmpty(t, log.UserName)
						assert.NotEmpty(t, log.Status)
						assert.NotEmpty(t, log.Feeling)
						assert.NotEmpty(t, log.Timestamp)
					}
				}
			}
		})
	}
}
