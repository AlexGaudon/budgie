package models

import (
	"database/sql"
	"fmt"
)

type CategoriesRepo struct {
	DB *sql.DB
}

func (r *CategoriesRepo) Find(userId string) ([]*Category, error) {
	query := `SELECT id, userid, name, created_at, updated_at, deleted_at FROM categories WHERE userid = $1 AND deleted_at IS NULL`

	rows, err := r.DB.Query(query, userId)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no categories found")
	}

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

	rows.Close()

	return categories, nil
}

func (r *CategoriesRepo) FindOne(c *Category) (*Category, error) {
	query := `SELECT id, userid, name, created_at, updated_at, deleted_at FROM categories WHERE id = $1 AND deleted_at IS NULL`

	row := r.DB.QueryRow(query, c.ID)

	err := row.Scan(
		&c.ID,
		&c.UserID,
		&c.Name,
		&c.CreatedAt,
		&c.UpdatedAt,
		&c.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("category with id not found")
	}

	if err != nil {
		return nil, err
	}

	return c, nil
}

func (r *CategoriesRepo) Exists(c *Category) bool {
	f, err := r.FindOne(c)

	if err != nil {
		return false
	}

	return f != nil && f.ID != ""
}

func (r *CategoriesRepo) Save(c *Category) (*Category, error) {
	return r.create(c)
}

func (r *CategoriesRepo) Delete(id string) error {
	query := `UPDATE categories SET deleted_at = (NOW() AT TIME ZONE 'UTC') WHERE id = $1`

	rows, err := r.DB.Query(query, id)

	rows.Close()

	if err != nil {
		return err
	}

	return nil
}

func (r *CategoriesRepo) create(c *Category) (*Category, error) {
	query := `INSERT INTO categories (userid, name)
	VALUES ($1, $2) RETURNING id, created_at, updated_at, deleted_at`

	row := r.DB.QueryRow(query, c.UserID, c.Name)

	err := row.Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt, &c.DeletedAt)

	if err != nil {
		return nil, err
	}

	if c.ID == "" {
		return nil, fmt.Errorf("error creating category")
	}

	return c, nil
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
