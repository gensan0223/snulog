package usecase

import (
	"context"

	"github.com/gensan0223/snulog/proto"
)

type LogUsecase interface {
	AddLogs(ctx context.Context, entry *proto.LogEntry) (*proto.AddResponse, error)
	FetchLogs(ctx context.Context) (*proto.FetchResponse, error)
}

type logUsecase struct {
	logs []*proto.LogEntry
}

func NewLogUsecase() LogUsecase {
	return &logUsecase{
		logs: []*proto.LogEntry{},
	}
}

func (u *logUsecase) AddLogs(ctx context.Context, entry *proto.LogEntry) (*proto.AddResponse, error) {
	u.logs = append(u.logs, entry)
	return &proto.AddResponse{Message: "added successfully"}, nil
}

func (u *logUsecase) FetchLogs(ctx context.Context) (*proto.FetchResponse, error) {
	return &proto.FetchResponse{
		Logs: u.logs,
	}, nil
}
