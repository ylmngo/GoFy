package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"gofy/internal/validator"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = errors.New("Duplicate Email")
	ErrRecordNotFound = errors.New("Record Not Found")
)

var AnonymousUser = &User{Username: "Anon"}

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Username  string    `json:"username"`
	Password  password  `json:"-"`
	Email     string    `json:"email"`
	Version   int       `json:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

type UserModel struct {
	DB *sql.DB
}

func (user *User) IsAnonymous() bool {
	return user == AnonymousUser
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return nil
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, err
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePassword(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must not be empty")
	v.Check(len(password) >= 8, "password", "must be greater than 8 bytes")
	v.Check(len(password) <= 72, "password", "must be less than 72 bytes")
}

func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Username != "nil", "username", "must be provided")
	v.Check(len(user.Username) <= 500, "username", "must be less than 500 bytes")

	ValidateEmail(v, user.Email)

	if user.Password.plaintext != nil {
		ValidatePassword(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

func (um UserModel) Insert(user *User) error {
	query := `
	INSERT INTO users (name, email, password_hash) 
	VALUES ($1, $2, $3)
	RETURNING id, created_at, version;`

	args := []interface{}{user.Username, user.Email, user.Password.hash}

	err := um.DB.QueryRow(query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (um UserModel) GetByEmail(email string) (*User, error) {
	var user User
	query := `
		SELECT id, created_at, name, email, password_hash, version 
		FROM users 
		WHERE email = $1`

	args := email

	err := um.DB.QueryRow(query, args).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

// Add Update and Delete
func (um UserModel) GetByToken(tokenScope string, tokenPlainText string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlainText))
	var user User
	query := `
	SELECT users.id, users.created_at, users.name, users.email, users.password_hash, users.version
	FROM users
	INNER JOIN tokens
	ON users.id = tokens.user_id
	WHERE tokens.hash = $1
	AND tokens.scope = $2
	AND tokens.expiry > $3`

	args := []interface{}{tokenHash[:], tokenScope, time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := um.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Username,
		&user.Email,
		&user.Password.hash,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
