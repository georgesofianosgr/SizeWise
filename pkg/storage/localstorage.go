package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/georgesofianosgr/sizewise/pkg/utils/fileutils"
)

var _ Storage = (*Local)(nil)

type Local struct {
	Base
	BasePath string `json:"basePath"`
}

func NewLocal(id, basePath string) Local {
	path := basePath
	endsWithSlash := strings.HasSuffix(basePath, "/")
	if !endsWithSlash {
		path = basePath + "/"
	}
	return Local{
		Base: Base{
			ID:   id,
			Type: "local",
		},
		BasePath: path,
	}
}

func (l Local) GetID() string {
	return l.ID
}

func (l Local) GetType() string {
	return l.Type
}

func (l Local) ShouldCache() bool {
	return l.Cache
}

func (l Local) EntryExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("error while checking local entry: %w", err)
	}

	return true, nil
}

func (l Local) WriteEntry(path string, reader io.Reader) error {
	filePath := filepath.Join(l.BasePath, path)
	if fileutils.FileExists(filePath) {
		return ErrEntryExists
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error while writing entry: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("error while writing entry: %w", err)
	}

	return nil
}

func (l Local) ReadEntry(path string, writer io.Writer) error {
	filePath := filepath.Join(l.BasePath, path)
	exists := fileutils.FileExists(filePath)
	if !exists {
		return ErrEntryNotFound
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error while writing entry: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(writer, file)
	return err
}
