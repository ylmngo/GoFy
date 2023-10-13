package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"gofy/internal/validator"
	"time"
)

const (
	ScopeAuthentication = "authentication"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expriy"`
	Scope     string    `json:"-"`
}

type TokenModel struct {
	DB *sql.DB
}

func ValidatePlainTextToken(v *validator.Validator, tokenPlainText string) {
	v.Check(len(tokenPlainText) != 0, "token", "must be provided")
	v.Check(len(tokenPlainText) == 26, "token", "must be 26 bytes long")
}

func (tm TokenModel) GenerateToken(userId int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userId,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	rndBytes := make([]byte, 16)
	// rand.Read() writes to rndBytes a Cryptographically Secure Pseudo Random Number through the OS's CSPRNG
	_, err := rand.Read(rndBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(rndBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	// convert hash which is a byte array to tok.Hash which is a byte slice
	token.Hash = hash[:]

	return token, nil
}

func (tm TokenModel) Insert(token *Token) error {
	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)`

	args := []interface{}{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := tm.DB.ExecContext(ctx, query, args...)

	return err
}
