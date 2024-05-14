package data

import (
	"context"
	"database/sql"
	"time"
)

type File struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	Size      int32     `json:"size"`
	Metadata  string    `json:"metadata"`
	Version   int32     `json:"-"`
}

type FileModel struct {
	DB *sql.DB
}

func (fm FileModel) Insert(f *File, userID int64) error {
	query := `
		INSERT into files (filename, metadata)
		VALUES ($1, $2)
		RETURNING uploaded_at;`

	args := []interface{}{f.Name, f.Metadata}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := fm.DB.QueryRowContext(ctx, query, args...).Scan(&f.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (fm FileModel) Get(fileId int64) (*File, error) {
	var f File
	query := `
	SELECT filename, metadata
	FROM files 
	WHERE file_id = $1`

	args := []interface{}{fileId}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := fm.DB.QueryRowContext(ctx, query, args...).Scan(
		&f.Name,
		&f.Metadata,
	)

	if err != nil {
		return nil, err
	}

	return &f, nil
}

func (fm FileModel) GetAll() ([]File, error) {
	var fs []File
	query := `
		SELECT filename, metadata
		FROM files 
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := fm.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var name, meta string
	for rows.Next() {
		if err = rows.Scan(&name, &meta); err != nil {
			return nil, err
		}
		fs = append(fs, File{Name: name, Metadata: meta})
	}

	return fs, nil
}
