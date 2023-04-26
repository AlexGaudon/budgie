package storage

import "time"

type Transaction struct {
	ID           string     `json:"id"`
	UserId       string     `json:"userid"`
	Vendor       string     `json:"vendor"`
	Description  string     `json:"description"`
	CategoryID   string     `json:"category_id"`
	CategoryName string     `json:"category_name"`
	Amount       int        `json:"amount"`
	Date         time.Time  `json:"date"`
	Type         string     `json:"type"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}
