package sizewise

import (
	"bytes"
	"fmt"
	"time"

	"github.com/georgesofianosgr/sizewise/pkg/config"
	"github.com/georgesofianosgr/sizewise/pkg/resizer"
	"github.com/georgesofianosgr/sizewise/pkg/storage"
	"golang.org/x/sync/errgroup"
)

var ErrResizeTimeout = fmt.Errorf("timeout while resizing image")

func ModifyAndCacheImage(config *config.Config, storageKey string, modifier string, path string, onEntryRetrieved func([]byte) error) error {
	modifiers, err := resizer.ParseModifiers(modifier)
	if err != nil {
		return err
	}

	foundStorage := storage.FindStorage(config.Storages, storageKey)
	if foundStorage == nil {
		return storage.ErrStorageNotFound
	}

	pathWithModifier := resizer.ModifiedPath(path, modifiers)
	modifiedExists, err := foundStorage.EntryExists(pathWithModifier)
	if err != nil {
		return err
	}

	if modifiedExists {
		content, err := readImage(foundStorage, pathWithModifier)
		if err != nil {
			return err
		}
		err = onEntryRetrieved(content)
		if err != nil {
			return err
		}
		return nil
	}

	// Resize pass content to callback and then write to storage
	content, err := resizeImage(foundStorage, path, modifiers)
	if err != nil {
		return err
	}

	group := errgroup.Group{}
	group.Go(func() error {
		return onEntryRetrieved(content)
	})

	if foundStorage.ShouldCache() {
		group.Go(func() error {
			return foundStorage.WriteEntry(pathWithModifier, bytes.NewReader(content))
		})
	}
	if err := group.Wait(); err != nil {
		return err
	}

	return nil
}

func readImage(s storage.Storage, path string) ([]byte, error) {
	content := bytes.Buffer{}
	err := s.ReadEntry(path, &content)
	if err != nil {
		return nil, err
	}

	return content.Bytes(), nil
}

func resizeImage(s storage.Storage, originalPath string, modifiers resizer.Modifiers) ([]byte, error) {
	content := bytes.Buffer{}
	err := s.ReadEntry(originalPath, &content)
	if err != nil {
		return nil, err
	}

	resizedFileBuffer := bytes.Buffer{}
	errChan := make(chan error, 1)

	// Avoid blocking the main thread
	go func() {
		errChan <- resizer.Resize(&content, &resizedFileBuffer, modifiers.CalculatedWidth, modifiers.CalculatedHeight)
	}()

	select {
	case err := <-errChan:
		if err != nil {
			return nil, err
		}
	case <-time.After(10 * time.Second):
		return nil, ErrResizeTimeout
	}

	return resizedFileBuffer.Bytes(), nil
}
