package data

import "time"

type File struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	Size      int32     `json:"size"`
	Version   int32     `json:"-"`
}
