package usecase

import (
	"context"

	"github.com/gensan0223/snulog/internal/repository"
	"github.com/gensan0223/snulog/proto"
)

type LogUsecase interface {
	AddLogs(ctx context.Context, entry *proto.LogEntry) (*proto.AddResponse, error)
	FetchLogs(ctx context.Context) (*proto.FetchResponse, error)
}

type logUsecase struct {
	repo repository.LogRepository
}

func NewLogUsecase(repo repository.LogRepository) LogUsecase {
	return &logUsecase{
		repo: repo,
	}
}

func (u *logUsecase) AddLogs(ctx context.Context, entry *proto.LogEntry) (*proto.AddResponse, error) {
	err := u.repo.Save(ctx, entry)
	if err != nil {
		return nil, err
	}
	return &proto.AddResponse{Message: "added successfully"}, nil
}

func (u *logUsecase) FetchLogs(ctx context.Context) (*proto.FetchResponse, error) {
	logs, err := u.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return &proto.FetchResponse{
		Logs: logs,
	}, nil
}
