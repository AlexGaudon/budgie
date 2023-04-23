package storage

import "time"

type Transaction struct {
	ID          string     `json:"id"`
	UserId      string     `json:"userid"`
	Description string     `json:"description"`
	Category    string     `json:"category"`
	Amount      int        `json:"amount"`
	Date        time.Time  `json:"date"`
	Type        string     `json:"type"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
