package storage

import (
	"io"
	"os"
)

func Write(path string, src io.Reader) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, src)
	return err
}

func Read(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	return data, err
}
