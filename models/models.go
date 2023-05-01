package models

import (
	"database/sql"
	"time"
)

type BaseModel struct {
	ID        string       `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`
}

type User struct {
	ID           string       `json:"id"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	DeletedAt    sql.NullTime `json:"deleted_at"`
	Username     string       `json:"username"`
	PasswordHash string       `json:"-"`
}

type Category struct {
	ID        string       `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"-"`
	UserID    string       `json:"user"`
	Name      string       `json:"name"`
}

type Budget struct {
	ID        string       `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"-"`

	UserID   string    `json:"user"`
	Name     string    `json:"name"`
	Category string    `json:"category"`
	Amount   int       `json:"amount"`
	Period   time.Time `json:"period"`
}
