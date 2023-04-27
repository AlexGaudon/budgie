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
		parent_id UUID REFERENCES category(id),
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
	query := `INSERT INTO category (name, userid, parent_id)
VALUES ($1, $2, $3) RETURNING id`

	fmt.Println("ParentID: ", c.ParentId)

	var id string
	var err error
	if c.ParentId == nil {
		// If ParentId is an empty string, set it to NULL in the query
		err = s.db.QueryRow(query, c.Name, c.UserID, nil).Scan(&id)
	} else {
		// Otherwise, pass ParentId as is
		err = s.db.QueryRow(query, c.Name, c.UserID, c.ParentId).Scan(&id)
	}

	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *PostgresStore) GetCategoryById(id string) (*Category, error) {
	query := `
SELECT
	c.*,
	p.name
AS parent_name
FROM category c
LEFT JOIN category p ON c.parent_id = p.id
WHERE
	c.deleted_at IS NULL 
	AND c.id = $1
	`

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
	rows, err := s.db.Query(`
SELECT
	c.*,
	p.name
AS parent_name
FROM category c
LEFT JOIN category p ON c.parent_id = p.id
WHERE c.deleted_at IS NULL AND c.userid = $1`, userId)

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

	var deletedAt sql.NullTime

	err := rows.Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&category.ParentId,
		&category.ParentName,
		&category.CreatedAt,
		&category.UpdatedAt,
		&deletedAt,
	)

	if deletedAt.Valid {
		category.DeletedAt = &deletedAt.Time
	}

	return category, err
}
