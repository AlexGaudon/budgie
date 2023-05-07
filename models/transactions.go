package models

import (
	"database/sql"
	"fmt"
	"time"
)

type TransactionsRepo struct {
	DB *sql.DB
}

type transactionPredicateFunction func(*Transaction) bool

func (r *TransactionsRepo) Filter(transactions []*Transaction, pred transactionPredicateFunction) []*Transaction {
	filtered := []*Transaction{}

	for _, transaction := range transactions {
		if pred(transaction) {
			filtered = append(filtered, transaction)
		}
	}

	return filtered
}

func (r *TransactionsRepo) Find(userId string) ([]*Transaction, error) {
	query := `
SELECT
	t.id,
	t.userid,
	t.amount,
	categories.NAME AS category_name,
	t.category,
	t.description,
	t.vendor,
	t.date,
	t.type,
	t.created_at,
	t.updated_at,
	t.deleted_at
FROM transactions t
JOIN categories
ON categories.id = t.category
WHERE t.deleted_at IS NULL
AND t.userid = $1
ORDER BY t.created_at DESC`

	rows, err := r.DB.Query(query, userId)

	if err != nil {
		return nil, err
	}

	transactions := []*Transaction{}

	for rows.Next() {
		transaction, err := scanIntoTransaction(rows)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	rows.Close()

	return transactions, nil
}

func (r *TransactionsRepo) FindOne(t *Transaction) (*Transaction, error) {
	query := `
SELECT t.id,
	t.userid,
	t.amount,
	categories.NAME AS category_name,
	t.category,
	t.description,
	t.vendor,
	t.date,
	t.type,
	t.created_at,
	t.updated_at,
	t.deleted_at
FROM transactions t
JOIN categories
ON categories.id = t.category
WHERE t.deleted_at IS NULL
AND t.id = $1`

	if t.ID == "" {
		return nil, fmt.Errorf("you must provide an id")
	}

	row := r.DB.QueryRow(query, t.ID)

	transaction := &Transaction{}

	err := row.Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.Amount,
		&transaction.Category,
		&transaction.CategoryID,
		&transaction.Description,
		&transaction.Vendor,
		&transaction.Date,
		&transaction.Type,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
		&transaction.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return transaction, nil
}

func (r *TransactionsRepo) Exists(t *Transaction) bool {
	_, err := r.FindOne(t)
	return err == nil
}

func (r *TransactionsRepo) Delete(id string) error {
	query := `UPDATE transactions SET deleted_at = (NOW() AT TIME ZONE 'UTC') WHERE id = $1`

	_, err := r.DB.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}

func (r *TransactionsRepo) Save(t *Transaction) (*Transaction, error) {
	if r.Exists(t) {
		return r.update(t)
	}
	return r.create(t)
}

func (r *TransactionsRepo) create(t *Transaction) (*Transaction, error) {
	query := `INSERT INTO transactions (userid, amount, category, description, vendor, date, type)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at, deleted_at`

	row := r.DB.QueryRow(query, t.UserID, t.Amount, t.CategoryID, t.Description.String, t.Vendor, t.Date, t.Type)

	err := row.Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt, &t.DeletedAt)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *TransactionsRepo) update(t *Transaction) (*Transaction, error) {
	query := `UPDATE transactions SET
	amount = $1,
	category = $2,
	description = $3,
	vendor = $4,
	date = $5,
	type = $6,
	updated_at = (NOW() AT TIME ZONE 'UTC')
	WHERE id = $7`

	_, err := r.DB.Exec(query, t.Amount, t.CategoryID, t.Description.String, t.Vendor, t.Date, t.Type, t.ID)
	if err != nil {
		return nil, err
	}

	t.UpdatedAt = time.Now().UTC()

	return t, nil
}

func scanIntoTransaction(rows *sql.Rows) (*Transaction, error) {
	t := &Transaction{}
	err := rows.Scan(
		&t.ID,
		&t.UserID,
		&t.Amount,
		&t.Category,
		&t.CategoryID,
		&t.Description,
		&t.Vendor,
		&t.Date,
		&t.Type,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.DeletedAt,
	)

	return t, err
}
