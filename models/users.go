package models

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func (u *User) IsPasswordValid(candidate string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(candidate)) == nil
}

func NewUser(username, password string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	return &User{
		Username:     username,
		PasswordHash: string(hash),
	}, nil
}

type UserRepo struct {
	DB *sql.DB
}

func (r *UserRepo) Save(u *User) (*User, error) {
	if r.Exists(u) {
		return r.update(u)
	}
	return r.create(u)
}

func (r *UserRepo) Exists(user *User) bool {
	f, err := r.FindOne(user)

	if err != nil {
		return false
	}

	return f != nil && f.ID != ""
}

func (r *UserRepo) FindOne(user *User) (*User, error) {
	query := ""
	paramOne := ""
	if user.ID != "" { // if the ID is provided, we use that to find it.
		query = `SELECT id, username, passwordhash, created_at, updated_at, deleted_at FROM users WHERE id = $1`
		paramOne = user.ID
	} else if user.Username != "" {
		query = `SELECT id, username, passwordhash, created_at, updated_at, deleted_at FROM users WHERE username = $1`
		paramOne = user.Username
	}

	row := r.DB.QueryRow(query, paramOne)

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return user, nil
}

func (r *UserRepo) create(u *User) (*User, error) {
	query := `
	INSERT INTO users (username, passwordhash) VALUES($1, $2) RETURNING id, created_at, updated_at, deleted_at`

	us := BaseModel{}

	row := r.DB.QueryRow(query, u.Username, u.PasswordHash)

	err := row.Scan(&us.ID, &us.CreatedAt, &us.UpdatedAt, &us.DeletedAt)

	if err != nil {
		return nil, err
	}

	if us.ID == "" {
		return nil, fmt.Errorf("a user with this name already exists")
	}

	return u, nil
}

func (r *UserRepo) update(u *User) (*User, error) {
	query := `UPDATE users SET username = $1, passwordhash = $2, updated_at = $3 WHERE id = $4 RETURNING updated_at`

	us := BaseModel{}

	err := r.DB.QueryRow(query, u.Username, u.PasswordHash, time.Now().UTC(), u.ID).Scan(&us.UpdatedAt)

	if err != nil {
		return nil, err
	}

	u.UpdatedAt = us.UpdatedAt

	return u, nil
}
