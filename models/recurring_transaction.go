package models

import (
	"database/sql"
	"fmt"
	"time"
)

type RecurringTransactionRepo struct {
	DB *sql.DB
}

func (r *RecurringTransactionRepo) Find(userId string) ([]*RecurringTransaction, error) {
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
	t.last_execution,
	t.next_execution,
	t.unit_of_measure,
	t.frequency_count,
	t.created_at,
	t.updated_at,
	t.deleted_at
FROM recurring_transactions t
JOIN categories
ON categories.id = t.category
WHERE t.deleted_at IS NULL
AND t.userid = $1
ORDER BY t.created_at DESC`

	rows, err := r.DB.Query(query, userId)

	if err != nil {
		return nil, err
	}

	recurringTransactions := []*RecurringTransaction{}

	for rows.Next() {
		rt, err := scanIntoRecurringTransaction(rows)
		if err != nil {
			return nil, err
		}

		recurringTransactions = append(recurringTransactions, rt)
	}

	rows.Close()

	return recurringTransactions, nil
}

func (r *RecurringTransactionRepo) FindAll_Internal() ([]*RecurringTransaction, error) {
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
	t.last_execution,
	t.next_execution,
	t.unit_of_measure,
	t.frequency_count,
	t.created_at,
	t.updated_at,
	t.deleted_at
FROM recurring_transactions t
JOIN categories
ON categories.id = t.category
WHERE t.deleted_at IS NULL
ORDER BY t.created_at DESC`

	rows, err := r.DB.Query(query)

	if err != nil {
		return nil, err
	}

	recurringTransactions := []*RecurringTransaction{}

	for rows.Next() {
		rt, err := scanIntoRecurringTransaction(rows)
		if err != nil {
			return nil, err
		}

		recurringTransactions = append(recurringTransactions, rt)
	}

	rows.Close()

	return recurringTransactions, nil
}

func (r *RecurringTransactionRepo) FindOne(t *RecurringTransaction) (*RecurringTransaction, error) {
	fmt.Println("FINDONE: ", t.ID)
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
	t.last_execution,
	t.next_execution,
	t.unit_of_measure,
	t.frequency_count,
	t.created_at,
	t.updated_at,
	t.deleted_at
FROM recurring_transactions t
JOIN categories
ON categories.id = t.category
WHERE t.deleted_at IS NULL
AND t.id = $1
ORDER BY t.created_at DESC`

	if t.ID == "" {
		return nil, fmt.Errorf("you must provide an id")
	}

	row := r.DB.QueryRow(query, t.ID)

	transaction := &RecurringTransaction{}

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
		&transaction.LastExecution,
		&transaction.NextExecution,
		&transaction.UnitOfMeasure,
		&transaction.Frequency,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
		&transaction.DeletedAt,
	)

	if err != nil {
		fmt.Println("ERR", err.Error())
		return nil, err
	}

	return transaction, nil
}

func (r *RecurringTransactionRepo) Exists(t *RecurringTransaction) bool {
	_, err := r.FindOne(t)
	return err == nil
}

func (r *RecurringTransactionRepo) Delete(id string) error {
	query := `UPDATE recurring_transactions SET deleted_at = (NOW() AT TIME ZONE 'UTC') WHERE id = $1`

	_, err := r.DB.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}

func (r *RecurringTransactionRepo) Save(t *RecurringTransaction) (*RecurringTransaction, error) {
	if r.Exists(t) {
		return r.update(t)
	}
	return r.create(t)
}

func (r *RecurringTransactionRepo) create(t *RecurringTransaction) (*RecurringTransaction, error) {
	query := `INSERT INTO recurring_transactions (userid, amount, category, description, vendor, date, type, last_execution, next_execution, unit_of_measure, frequency_count)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id, created_at, updated_at, deleted_at`

	row := r.DB.QueryRow(query, t.UserID, t.Amount, t.CategoryID, t.Description.String, t.Vendor, t.Date, t.Type, t.LastExecution, t.NextExecution, t.UnitOfMeasure, t.Frequency)

	err := row.Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt, &t.DeletedAt)

	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *RecurringTransactionRepo) update(t *RecurringTransaction) (*RecurringTransaction, error) {
	fmt.Println("UPDATING", t.ID)
	query := `UPDATE recurring_transactions SET
	amount = $1,
	category = $2,
	description = $3,
	vendor = $4,
	date = $5,
	type = $6,
	last_execution = $7,
	next_execution = $8,
	unit_of_measure = $9,
	frequency_count = $10,
	updated_at = (NOW() AT TIME ZONE 'UTC')
	WHERE id = $11`

	_, err := r.DB.Exec(query,
		t.Amount,
		t.CategoryID,
		t.Description.String,
		t.Vendor,
		t.Date,
		t.Type,
		t.LastExecution,
		t.NextExecution,
		t.UnitOfMeasure,
		t.Frequency,
		t.ID,
	)

	if err != nil {
		return nil, err
	}

	t.UpdatedAt = time.Now().UTC()

	return t, nil
}

func scanIntoRecurringTransaction(rows *sql.Rows) (*RecurringTransaction, error) {
	t := &RecurringTransaction{}
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
		&t.LastExecution,
		&t.NextExecution,
		&t.UnitOfMeasure,
		&t.Frequency,
		&t.CreatedAt,
		&t.UpdatedAt,
		&t.DeletedAt,
	)

	return t, err
}
