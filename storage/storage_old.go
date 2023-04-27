package storage

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/alexgaudon/budgie/config"
	_ "github.com/lib/pq"
)

var DB *PostgresStore

type PostgresStore struct {
	db *sql.DB
}

func (s *PostgresStore) enableUUID() error {
	query := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`

	_, err := s.db.Exec(query)

	if err != nil {
		log.Println("ERROR ON QUERY: ", query)
		return err
	}

	return nil
}

func SetupDatabase() error {
	config := config.GetConfig()
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBName, config.DBUserName, config.DBUserPassword)
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	log.Println("Connected to database")

	DB = &PostgresStore{
		db: db,
	}

	return nil
}

func (s *PostgresStore) Init() error {
	err := s.enableUUID()
	if err != nil {
		return err
	}

	err = s.createCategoryTable()
	if err != nil {
		return err
	}

	err = s.createTransactionTable()
	if err != nil {
		return err
	}

	err = s.createUserTable()
	if err != nil {
		return err
	}

	err = s.createBudgetTable()
	if err != nil {
		return err
	}

	return nil
}
