package storage

import (
	"io"
	"os"
)

func Write(path string, data io.Reader) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, data)
	return err
}

func Read(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	return data, err
}
