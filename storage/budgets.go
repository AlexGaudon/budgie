package storage

import (
	"database/sql"
	"fmt"
	"log"
)

func (s *PostgresStore) createBudgetTable() error {
	query := `CREATE TABLE IF NOT EXISTS budgets (
		id UUID PRIMARY KEY DEFAULT public.uuid_generate_v4(),
		userid UUID NOT NULL,
		category VARCHAR(255) NOT NULL,
		amount INTEGER NOT NULL,
		recurring BOOLEAN NOT NULL DEFAULT FALSE,
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

func (s *PostgresStore) CreateBudget(b *Budget) (string, error) {
	query := `INSERT INTO budgets (userid, category, amount, recurring)
	VALUES ($1, $2, $3, $4) RETURNING id`

	var id string
	err := s.db.QueryRow(query, b.UserId, b.Category, b.Amount, b.Recurring).Scan(&id)

	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *PostgresStore) GetBudgets(userId string) ([]*Budget, error) {
	rows, err := s.db.Query("SELECT * FROM budgets WHERE deleted_at IS NULL AND userid = $1", userId)

	if err != nil {
		return nil, err
	}

	budgets := []*Budget{}

	defer rows.Close()

	for rows.Next() {
		budget, err := scanIntoBudget(rows)
		if err != nil {
			return nil, err
		}
		budgets = append(budgets, budget)
	}

	return budgets, nil
}

func (s *PostgresStore) GetBudgetById(id string) (*Budget, error) {
	rows, err := s.db.Query("SELECT * FROM budgets WHERE id = $1 AND deleted_at IS NULL", id)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoBudget(rows)
	}

	return nil, fmt.Errorf("budget %s not found", id)
}

func (s *PostgresStore) DeleteBudget(id string, userId string) error {
	_, err := s.GetBudgetById(id)

	if err != nil {
		return err
	}

	_, err = s.db.Query("UPDATE budgets SET deleted_at = (NOW() AT TIME ZONE 'UTC') WHERE id = $1 AND userid = $2", id, userId)

	return err
}

func scanIntoBudget(rows *sql.Rows) (*Budget, error) {
	budget := new(Budget)
	err := rows.Scan(
		&budget.ID,
		&budget.UserId,
		&budget.Category,
		&budget.Amount,
		&budget.Recurring,
		&budget.CreatedAt,
		&budget.UpdatedAt,
		&budget.DeletedAt,
	)

	return budget, err
}
