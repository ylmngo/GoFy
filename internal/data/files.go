package data

import (
	"context"
	"database/sql"
	"time"
)

type File struct {
	Filename   string    `json:"file_name"`
	FileID     []byte    `json:"file_id"`
	Metadata   string    `json:"metadata"`
	UserID     int64     `json:"user_id"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type FileModel struct {
	DB *sql.DB
}

func (fm FileModel) Insert(f *File, userID int64) error {
	query := `
		INSERT into files (user_id, filename, file_id, metadata)
		VALUES ($1, $2, $3, $4)
		RETURNING uploaded_at;`

	args := []interface{}{userID, f.Filename, f.FileID, f.Metadata}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := fm.DB.QueryRowContext(ctx, query, args...).Scan(&f.UploadedAt)
	if err != nil {
		return err
	}

	return nil
}

func (fm FileModel) GetFiles(userID int64) ([]File, error) {
	var files []File
	query := `
		SELECT user_id, filename, file_id, metadata, uploaded_at
		FROM files 
		WHERE user_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	rows, err := fm.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var f File
		if err := rows.Scan(&f.UserID, &f.Filename, &f.FileID, &f.Metadata, &f.UploadedAt); err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, nil
}

func (fm FileModel) GetFile(userID int64, filename string) (*File, error) {
	var f File
	query := `
	SELECT user_id, filename, file_id, metadata
	FROM files 
	WHERE user_id = $1
	AND filename = $2`

	args := []interface{}{userID, filename}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := fm.DB.QueryRowContext(ctx, query, args...).Scan(
		&f.UserID,
		&f.Filename,
		&f.FileID,
		&f.Metadata,
	)

	if err != nil {
		return nil, err
	}

	return &f, nil
}
