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
		description VARCHAR(255) NOT NULL,
		category VARCHAR(255) NOT NULL,
		amount INTEGER NOT NULL,
		date TIMESTAMP NOT NULL,
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

func (s *PostgresStore) CreateTransaction(t *Transaction) (string, error) {
	query := `INSERT INTO transactions (userid, description, category, amount, date)
	VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var id string
	err := s.db.QueryRow(query, t.UserId, t.Description, t.Category, t.Amount, t.Date).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *PostgresStore) GetTransactions(userId string) ([]*Transaction, error) {
	rows, err := s.db.Query("SELECT * FROM transactions WHERE deleted_at IS NULL AND userid = $1", userId)

	if err != nil {
		return nil, err
	}

	transactions := []*Transaction{}

	defer rows.Close()

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
	rows, err := s.db.Query("SELECT * FROM transactions WHERE id = $1 AND deleted_at IS NULL", id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoTransaction(rows)
	}

	return nil, fmt.Errorf("transaction %s not found", id)
}

func (s *PostgresStore) UpdateTransaction(t *Transaction) error { return nil }

func (s *PostgresStore) DeleteTransaction(id string, userId string) error {
	_, err := s.GetTransactionById(id)

	if err != nil {
		return err
	}

	_, err = s.db.Query("UPDATE transactions SET deleted_at = (NOW() AT TIME ZONE 'UTC') WHERE id = $1 AND userid = $2", id, userId)

	return err
}

func scanIntoTransaction(rows *sql.Rows) (*Transaction, error) {
	transaction := new(Transaction)
	err := rows.Scan(
		&transaction.ID,
		&transaction.UserId,
		&transaction.Description,
		&transaction.Category,
		&transaction.Amount,
		&transaction.Date,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
		&transaction.DeletedAt,
	)

	return transaction, err
}
