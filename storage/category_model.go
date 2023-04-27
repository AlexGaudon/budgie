package storage

import "time"

type Category struct {
	ID         string     `json:"id"`
	UserID     string     `json:"userid"`
	Name       string     `json:"name"`
	ParentId   *string    `json:"parent_id"`
	ParentName string     `json:"parent_name"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}
