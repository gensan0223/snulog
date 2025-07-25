package repository

import (
	"context"

	"github.com/gensan0223/snulog/proto"
)

type InMemoryLogRepository struct {
	logs []*proto.LogEntry
}

func NewInMemoryLogRepository() *InMemoryLogRepository {
	return &InMemoryLogRepository{
		logs: []*proto.LogEntry{},
	}
}

func (r *InMemoryLogRepository) Save(ctx context.Context, entry *proto.LogEntry) error {
	r.logs = append(r.logs, entry)
	return nil
}

func (r *InMemoryLogRepository) FindAll(ctx context.Context) ([]*proto.LogEntry, error) {
	return r.logs, nil
}
