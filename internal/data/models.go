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
}

func NewModels(db *sql.DB) Models {
	return Models{
		Files: FileModel{DB: db},
	}
}
