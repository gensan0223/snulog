package repository

import (
	"database/sql"
)

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
}

type UserRepository interface {
	GetUserByUsername(username string) (*User, error)
	CreateUser(username, passwordHash string) error
}

type postgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) GetUserByUsername(username string) (*User, error) {
	user := &User{}
	query := "SELECT id, username, password_hash FROM users WHERE username = $1"

	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *postgresUserRepository) CreateUser(username, passwordHash string) error {
	query := "INSERT INTO users (username, password_hash) VALUES ($1, $2)"
	_, err := r.db.Exec(query, username, passwordHash)
	return err
}
