package storage

import (
	"database/sql"
	"fmt"
	"log"
)

func (s *PostgresStore) createUserTable() error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT public.uuid_generate_v4(),
		username varchar(255) NOT NULL,
		passwordhash varchar(1024) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
		updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
		deleted_at TIMESTAMP
	)`

	_, err := s.db.Exec(query)

	if err != nil {
		log.Println("ERROR ON QUERY: ", query)
		return err
	}

	return nil
}

func (s *PostgresStore) CreateUser(u *User) (string, error) {
	query := `INSERT INTO users (username, passwordhash)
	VALUES ($1, $2) returning id`

	var id string
	err := s.db.QueryRow(query, u.Username, u.PasswordHash).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *PostgresStore) GetUserById(id string) (*User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	rows, err := s.db.Query(query, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, fmt.Errorf("user %s not found", id)
}

func (s *PostgresStore) GetUserByUsername(username string) (*User, error) {
	query := `SELECT * FROM users WHERE username = $1`
	rows, err := s.db.Query(query, username)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, fmt.Errorf("user %s not found", username)
}

func scanIntoUser(rows *sql.Rows) (*User, error) {
	user := new(User)

	err := rows.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	return user, err
}
