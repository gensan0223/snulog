package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/gensan0223/snulog/proto"
	"github.com/stretchr/testify/assert"
)

func TestAddLogsAndFetchLogs(t *testing.T) {
	uc := NewLogUsecase()

	entry := &proto.LogEntry{
		UserName:  "tester",
		Status:    "good",
		Feeling:   "🆒",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	_, err := uc.AddLogs(context.Background(), entry)
	assert.NoError(t, err)

	res, err := uc.FetchLogs(context.Background())
	assert.NoError(t, err)
	assert.Len(t, res.Logs, 1)
	assert.Equal(t, "tester", res.Logs[0].UserName)
}
