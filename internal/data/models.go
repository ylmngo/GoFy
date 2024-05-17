package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Files FileModel
	Users UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Files: FileModel{DB: db},
		Users: UserModel{DB: db},
	}
}
