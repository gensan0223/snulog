package repository

import (
	"github.com/gensan0223/snulog/proto"
	"golang.org/x/net/context"
)

type LogRepository interface {
	Save(ctx context.Context, entry *proto.LogEntry) error
	FindAll(ctx context.Context) ([]*proto.LogEntry, error)
}
