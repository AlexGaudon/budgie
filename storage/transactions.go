package storage

import (
	"database/sql"
	"fmt"
	"log"
)

func (s *PostgresStore) createTransactionTable() error {
	query := `CREATE TABLE IF NOT EXISTS transactions (
		id UUID PRIMARY KEY DEFAULT public.uuid_generate_v4(),
		userid UUID NOT NULL,
		vendor VARCHAR(255) NOT NULL,
		description VARCHAR(255) NOT NULL,
		category_id UUID NOT NULL REFERENCES category(id),
		amount INTEGER NOT NULL,
		date TIMESTAMP NOT NULL,
		type VARCHAR(255) NOT NULL,
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

func (s *PostgresStore) CreateTransaction(t *Transaction) (*Transaction, error) {
	query := `INSERT INTO transactions (userid, vendor, description, category_id, amount, date, type)
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	var id string
	err := s.db.QueryRow(query, t.UserId, t.Vendor, t.Description, t.CategoryID, t.Amount, t.Date, t.Type).Scan(&id)
	if err != nil {
		return &Transaction{}, err
	}

	return s.GetTransactionById(id)
}

func (s *PostgresStore) GetTransactionsByCategory(userId, category string) ([]*Transaction, error) {
	query := `
SELECT
	t.id,
	t.userid,
	t.vendor,
	t.description,
	t.category_id,
	c.name as category_name,
	t.amount,
	t.date,
	t.type,
	t.created_at,
	t.updated_at,
	t.deleted_at
FROM
	transactions t
JOIN
	category c ON t.category_id = c.id
WHERE
	t.deleted_at IS NULL
	AND t.userid = $1
	AND t.category_id = $2
ORDER BY
	t.created_at DESC;`

	rows, err := s.db.Query(query, userId, category)

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

	return transactions, nil
}

func (s *PostgresStore) GetTransactions(userId string) ([]*Transaction, error) {
	rows, err := s.db.Query(`
SELECT 
    t.id,
    t.userid,
	t.vendor,
    t.description,
	t.category_id,
    c.name as category_name,
    t.amount,
    t.date,
    t.type,
    t.created_at,
    t.updated_at,
    t.deleted_at
FROM 
    transactions t
JOIN 
    category c ON t.category_id = c.id
WHERE 
    t.deleted_at IS NULL
    AND t.userid = $1
ORDER BY 
    t.created_at DESC;`,
		userId)

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

	return transactions, nil
}

func (s *PostgresStore) GetTransactionById(id string) (*Transaction, error) {
	rows, err := s.db.Query(`
SELECT 
	t.id,
    t.userid,
	t.vendor,
    t.description,
	t.category_id,
    c.name as category_name,
    t.amount,
    t.date,
    t.type,
    t.created_at,
    t.updated_at,
    t.deleted_at
FROM 
    transactions t
JOIN 
    category c ON t.category_id = c.id
WHERE 
    t.id = $1
    AND t.deleted_at IS NULL;
`, id)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		return scanIntoTransaction(rows)
	}

	return nil, fmt.Errorf("transaction %s not found", id)
}

func (s *PostgresStore) UpdateTransaction(t *Transaction) error {
	query := `UPDATE transactions SET vendor = $1, description=$2, category_id=$3, amount=$4, date=$5, type=$6, updated_at = (NOW() AT TIME ZONE 'UTC') WHERE id = $7`

	rows, err := s.db.Query(query, t.Vendor, t.Description, t.CategoryID, t.Amount, t.Date, t.Type, t.ID)

	if err != nil {
		return err
	}

	defer rows.Close()

	return nil
}

func (s *PostgresStore) DeleteTransaction(id string) error {
	_, err := s.GetTransactionById(id)

	if err != nil {
		return err
	}

	rows, err := s.db.Query("UPDATE transactions SET deleted_at = (NOW() AT TIME ZONE 'UTC') WHERE id = $1", id)

	rows.Close()

	if err != nil {
		return err
	}

	return nil
}

func scanIntoTransaction(rows *sql.Rows) (*Transaction, error) {
	transaction := &Transaction{}
	err := rows.Scan(
		&transaction.ID,
		&transaction.UserId,
		&transaction.Vendor,
		&transaction.Description,
		&transaction.CategoryID,
		&transaction.CategoryName,
		&transaction.Amount,
		&transaction.Date,
		&transaction.Type,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
		&transaction.DeletedAt,
	)

	return transaction, err
}
