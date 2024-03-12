package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/georgesofianosgr/sizewise/internal/sizewise"
	"github.com/georgesofianosgr/sizewise/pkg/config"
	"github.com/georgesofianosgr/sizewise/pkg/resizer"
	"github.com/georgesofianosgr/sizewise/pkg/storage"
)

var ErrResizeTimeout = fmt.Errorf("timeout while resizing image")

func requestHandler(config *config.Config) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		modifier := r.PathValue("modifier")
		storageKey := r.PathValue("storage")
		path := r.PathValue("path")

		writeResponse := func(content []byte) error {
			contentType := http.DetectContentType(content)
			w.Header().Set("Content-Type", contentType)
			_, err := w.Write(content)
			if err != nil {
				http.Error(w, "Error while writing response", http.StatusInternalServerError)
				return err
			}

			return nil
		}

		err := sizewise.ModifyAndCacheImage(config, storageKey, modifier, path, writeResponse)
		if errors.Is(err, storage.ErrEntryNotFound) {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		if errors.Is(err, resizer.ErrInvalidModifier) {
			http.Error(w, "modifier parse error", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, "error while checking file", http.StatusInternalServerError)
		}
	}
}

func Start(config *config.Config, port string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{storage}/{modifier}/{path...}", requestHandler(config))

	server := &http.Server{
		Addr:         "localhost:" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	slog.Info("Server is running", "port", port)
	err := server.ListenAndServe()
	return err
}
