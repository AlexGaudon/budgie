package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"

	"github.com/alexgaudon/budgie/config"
	"github.com/alexgaudon/budgie/models"
)

type DBStore struct {
	migrationPath string
	db            *sql.DB
	User          *models.UserRepo
	Categories    *models.CategoriesRepo
	Budgets       *models.BudgetsRepo
}

func (d *DBStore) Initialize() error {
	d.User = &models.UserRepo{
		DB: d.db,
	}

	d.Categories = &models.CategoriesRepo{
		DB: d.db,
	}

	d.Budgets = &models.BudgetsRepo{
		DB: d.db,
	}

	err := d.handleMigrations()

	return err
}

func (d *DBStore) hasExecutedMigration(id string) bool {
	query := `SELECT executed FROM migrations WHERE id = $1`

	var executed bool
	err := d.db.QueryRow(query, id).Scan(&executed)

	if err != nil {
		return false
	}

	return executed
}

func (d *DBStore) runMigration(migration string) error {
	_, err := d.db.Exec(migration)

	if err != nil {
		return err
	}

	return nil
}

func (d *DBStore) markMigration(id string) error {
	query := `INSERT INTO migrations (id, executed)
	VALUES($1, true)`

	_, err := d.db.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}

func (d *DBStore) handleMigrations() error {
	migrationTableQuery := `CREATE TABLE IF NOT EXISTS migrations (
		id SERIAL PRIMARY KEY,
		executed BOOLEAN NOT NULL DEFAULT FALSE
	);`

	_, err := d.db.Exec(migrationTableQuery)

	if err != nil {
		log.Println("Error creating Migrations Table, ", err.Error())
		return err
	}

	files, err := os.ReadDir(d.migrationPath)

	for _, e := range files {
		if !e.IsDir() {
			migration := e.Name()
			id := strings.Split(migration, ".")[0]

			if !d.hasExecutedMigration(id) {
				path := d.migrationPath + "/" + e.Name()
				content, err := os.ReadFile(path)

				if err != nil {
					log.Printf("ERROR: Failed reading migration (%s): %s", id, err.Error())
				}

				err = d.runMigration(string(content))

				if err != nil {
					log.Printf("ERROR: Failed running migration (%s): %s", id, err.Error())
				}

				err = d.markMigration(id)

				if err != nil {
					log.Printf("ERROR: Failed marking migration (%s): %s", id, err.Error())
				}

				log.Println("Ran and marked migration: ", id)
			}
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func ConnectDatabase(migrationPath string) (*DBStore, error) {
	config := config.GetConfig()
	connStr := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBName, config.DBUserName, config.DBUserPassword)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to database")

	return &DBStore{
		migrationPath: migrationPath,
		db:            db,
	}, nil
}
