package models

import (
	"database/sql"
	"encoding/json"
	"reflect"
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

	UserID     string    `json:"user"`
	Category   string    `json:"category"`
	CategoryID string    `json:"category_id"`
	Amount     int       `json:"amount"`
	Period     time.Time `json:"period"`
}

type Transaction struct {
	ID        string       `json:"id"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at"`

	UserID      string         `json:"user"`
	Amount      int            `json:"amount"`
	Category    string         `json:"category"`
	CategoryID  string         `json:"category_id"`
	Description OptionalString `json:"description"`
	Vendor      string         `json:"vendor"`
	Date        time.Time      `json:"date"`
	Type        string         `json:"type"`
}

type RecurringTransaction struct {
	Transaction

	LastExecution time.Time `json:"last_execution"`
	NextExecution time.Time `json:"next_execution"`
	UnitOfMeasure string    `json:"unit_of_measure"`
	Frequency     int       `json:"frequency_count"`
}

type OptionalString sql.NullString

func (os *OptionalString) Scan(value interface{}) error {
	var s sql.NullString

	if err := s.Scan(value); err != nil {
		return err
	}

	if reflect.TypeOf(value) == nil {
		*os = OptionalString{s.String, false}
	} else {
		*os = OptionalString{s.String, true}
	}
	return nil
}

func (os *OptionalString) MarshalJSON() ([]byte, error) {
	if !os.Valid {
		return []byte("null"), nil
	}

	return json.Marshal(os.String)
}

func (os *OptionalString) UnmarshalJSON(b []byte) error {
	err := json.Unmarshal(b, &os.String)
	os.Valid = (err == nil)
	return err
}
