package models

import "database/sql"

type TransationsRepo struct {
	DB *sql.DB
}

func (r *TransationsRepo) Find(userId string) ([]*Transaction, error) {

	return nil, nil
}

func (r *TransationsRepo) FindOne(transaction *Transaction) (*Transaction, error) {
	return nil, nil
}

func (r *TransationsRepo) Exists(t *Transaction) bool {
	f, err := r.FindOne(t)

	if err != nil {
		return false
	}

	return f != nil && f.ID != ""
}

func (r *TransationsRepo) Delete(id string) error {
	return nil
}

func (r *TransationsRepo) create(t *Transaction) (*Transaction, error) {
	return nil, nil
}

func (r *TransationsRepo) update(t *Transaction) (*Transaction, error) {
	panic("NOT IMPLEMENTED")
}

func scanIntoTransaction(rows *sql.Rows) (*Transaction, error) {
	t := &Transaction{}
	err := rows.Scan(
		&t.ID,
		&t.ID,
		&t.ID,
		&t.ID,
		&t.ID,
		&t.ID,
		&t.ID,
	)

	return t, err
}
