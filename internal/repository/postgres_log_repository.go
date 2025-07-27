package repository

import (
	"context"
	"database/sql"

	"github.com/gensan0223/snulog/proto"
)

type PostgresLogRepository struct {
	db *sql.DB
}

func NewPostgresLogRepository(db *sql.DB) *PostgresLogRepository {
	return &PostgresLogRepository{db: db}
}

func (r *PostgresLogRepository) Save(ctx context.Context, entry *proto.LogEntry) error {
	_, err := r.db.ExecContext(ctx, `
        INSERT INTO logs (user_name, status, feeling, timestamp)
        VALUES ($1, $2, $3, $4)
        `, entry.UserName, entry.Status, entry.Feeling, entry.Timestamp)
	return err
}

func (r *PostgresLogRepository) FindAll(ctx context.Context) ([]*proto.LogEntry, error) {
	rows, err := r.db.QueryContext(ctx, `
        SELECT user_name, status, feeling, timestamp FROM logs ORDER BY timestamp desc
        `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []*proto.LogEntry
	for rows.Next() {
		var entry proto.LogEntry
		if err := rows.Scan(&entry.UserName, &entry.Status, &entry.Feeling, &entry.Timestamp); err != nil {
			return nil, err
		}
		logs = append(logs, &entry)
	}
	return logs, nil
}
