package data

import (
	"context"
	"database/sql"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Version   int       `json:"-"`
}

func (um UserModel) Insert(user *User) error {
	query := `
	INSERT INTO users (name, email, password_hash) 
	VALUES ($1, $2, $3)
	RETURNING id, created_at, version;`

	args := []interface{}{user.Name, user.Email, user.Password.hash}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := um.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)

	return err
}

func (um UserModel) GetByEmail(email string) (*User, error) {
	var user User
	query := `
		SELECT id, created_at, name, email, password_hash, version 
		FROM users 
		WHERE email = $1`

	args := email

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := um.DB.QueryRowContext(ctx, query, args).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Version,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (um UserModel) GetById(id int) (*User, error) {
	var user User
	query := `
		SELECT id, created_at, name, email, password_hash, version
		FROM users 
		WHERE id = $1
	`

	args := id
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := um.DB.QueryRowContext(ctx, query, args).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Version,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
