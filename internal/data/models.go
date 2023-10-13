package data

import "database/sql"

type Model struct {
	Users  UserModel
	Files  FileModel
	Tokens TokenModel
}

func NewModel(db *sql.DB) Model {
	model := Model{
		Users:  UserModel{DB: db},
		Tokens: TokenModel{DB: db},
		Files:  FileModel{DB: db},
	}

	return model
}
