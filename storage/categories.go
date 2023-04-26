package storage

import (
	"database/sql"
	"fmt"
	"log"
)

func (s *PostgresStore) createCategoryTable() error {
	query := `CREATE TABLE IF NOT EXISTS category (
		id UUID PRIMARY KEY DEFAULT public.uuid_generate_v4(),
		userid UUID NOT NULL,
		name VARCHAR(255) NOT NULL UNIQUE, -- Add UNIQUE constraint on name column
		created_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
		updated_at TIMESTAMP NOT NULL DEFAULT (NOW() AT TIME ZONE 'UTC'),
		deleted_at TIMESTAMP
	);`

	_, err := s.db.Exec(query)

	if err != nil {
		log.Println("ERROR ON QUERY: ", query)
		return err
	}

	return nil
}

func (s *PostgresStore) CreateCategory(c *Category) (string, error) {
	query := `INSERT INTO category (name, userid)
VALUES ($1, $2) RETURNING id`

	var id string
	err := s.db.QueryRow(query, c.Name, c.UserID).Scan(&id)

	if err != nil {
		return "", err

	}

	return id, nil
}

func (s *PostgresStore) GetCategoryById(id string) (*Category, error) {
	query := `SELECT * FROM category WHERE deleted_at IS NULL and id = $1`

	rows, err := s.db.Query(query, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		return scanIntoCategory(rows)
	}

	return nil, fmt.Errorf("transaction %s not found", id)
}

func (s *PostgresStore) GetCategories(userId string) ([]*Category, error) {
	rows, err := s.db.Query("SELECT * FROM category WHERE deleted_at IS NULL AND userid = $1 ORDER BY created_at DESC", userId)

	if err != nil {
		return nil, err
	}

	categories := []*Category{}

	for rows.Next() {
		category, err := scanIntoCategory(rows)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func (s *PostgresStore) DeleteCategory(id string) error {
	_, err := s.GetCategoryById(id)

	if err != nil {
		return err
	}

	rows, err := s.db.Query(`UPDATE transactions SET deleted_at = (NOW() AT TIME ZONE 'UTC') WHERE id = $1`, id)

	rows.Close()

	if err != nil {
		return err
	}

	return nil
}

func scanIntoCategory(rows *sql.Rows) (*Category, error) {
	category := &Category{}
	err := rows.Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&category.CreatedAt,
		&category.UpdatedAt,
		&category.DeletedAt,
	)

	return category, err
}
