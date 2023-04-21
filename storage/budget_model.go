package storage

import (
	"time"
)

type Budget struct {
	ID        string     `json:"id"`
	UserId    string     `json:"userid"`
	Category  string     `json:"category"`
	Amount    int        `json:"amount"`
	Recurring bool       `json:"recurring"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
