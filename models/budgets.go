package models

import (
	"database/sql"
	"fmt"
)

type BudgetsRepo struct {
	DB *sql.DB
}

func (r *BudgetsRepo) Find(userId string) ([]*Budget, error) {
	query := `SELECT
	budgets.id,
	budgets.userid,
	categories.name as category_name,
	budgets.category,
	budgets.amount,
	budgets.period,
	budgets.created_at,
	budgets.updated_at,
	budgets.deleted_at 
FROM budgets 
JOIN categories ON categories.id = budgets.category 
WHERE budgets.deleted_at IS NULL
AND budgets.userid = $1`

	rows, err := r.DB.Query(query, userId)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no categories found")
	}

	if err != nil {
		return nil, err
	}

	budgets := []*Budget{}

	for rows.Next() {
		budget, err := scanIntoBudget(rows)
		if err != nil {
			return nil, err
		}
		budgets = append(budgets, budget)
	}

	rows.Close()

	return budgets, nil
}

func (r *BudgetsRepo) FindOne(b *Budget) (*Budget, error) {
	query := `SELECT
	budgets.id,
	budgets.userid,
	categories.name as category_name,
	budgets.category,
	budgets.amount,
	budgets.period,
	budgets.created_at,
	budgets.updated_at,
	budgets.deleted_at 
FROM budgets 
JOIN categories ON categories.id = budgets.category 
WHERE budgets.deleted_at IS NULL
AND budgets.id = $1`

	if b.ID == "" {
		return nil, fmt.Errorf("you must provide an id")
	}

	row := r.DB.QueryRow(query, b.ID)

	err := row.Scan(
		&b.ID,
		&b.UserID,
		&b.Category,
		&b.CategoryID,
		&b.Amount,
		&b.Period,
		&b.CreatedAt,
		&b.UpdatedAt,
		&b.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("budget with id not found")
	}

	if err != nil {
		return nil, err
	}

	return b, nil
}

func (r *BudgetsRepo) Exists(b *Budget) bool {
	f, err := r.FindOne(b)

	if err != nil {
		return false
	}

	return f != nil && f.ID != ""
}

func (r *BudgetsRepo) Save(b *Budget) (*Budget, error) {
	if r.Exists(b) {
		return r.update(b)
	}
	return r.create(b)
}

func (r *BudgetsRepo) Delete(id string) error {
	query := `UPDATE budgets SET deleted_at = (NOW() AT TIME ZONE 'UTC') WHERE id = $1`

	rows, err := r.DB.Query(query, id)

	rows.Close()

	if err != nil {
		return err
	}

	return nil
}

func (r *BudgetsRepo) create(b *Budget) (*Budget, error) {
	query := `INSERT INTO budgets (userid, category, amount, period)
	VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at, deleted_at`

	row := r.DB.QueryRow(query, b.UserID, b.Category, b.Amount, b.Period)

	err := row.Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt, &b.DeletedAt)

	if err != nil {
		return nil, err
	}

	if b.ID == "" {
		return nil, fmt.Errorf("error creating budget")
	}

	return b, nil
}

func (r *BudgetsRepo) update(b *Budget) (*Budget, error) {
	panic("NOT IMPLEMENTED")
}

func scanIntoBudget(rows *sql.Rows) (*Budget, error) {
	b := &Budget{}
	err := rows.Scan(
		&b.ID,
		&b.UserID,
		&b.Category,
		&b.CategoryID,
		&b.Amount,
		&b.Period,
		&b.CreatedAt,
		&b.UpdatedAt,
		&b.DeletedAt,
	)

	return b, err
}