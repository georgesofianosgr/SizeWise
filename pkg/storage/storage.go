package storage

import (
	"fmt"
	"io"
)

type Base struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Cache bool   `json:"cache"`
}

type Storage interface {
	GetID() string
	GetType() string
	ShouldCache() bool
	EntryExists(path string) (bool, error)
	WriteEntry(path string, reader io.Reader) error
	ReadEntry(path string, writer io.Writer) error
}

var (
	ErrStorageNotFound = fmt.Errorf("storage not found")
	ErrEntryExists     = fmt.Errorf("entry already exists")
	ErrEntryNotFound   = fmt.Errorf("entry not found")
)
