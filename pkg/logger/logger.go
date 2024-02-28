package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

var ErrInvalidLogType = fmt.Errorf("invalid log type, log type can be text or json")

func newLogger(logType string, w io.Writer) (*slog.Logger, error) {
	switch strings.ToLower(logType) {
	case "text", "":
		return slog.New(slog.NewTextHandler(w, nil)), nil
	case "json":
		return slog.New(slog.NewJSONHandler(w, nil)), nil
	default:
		return nil, ErrInvalidLogType
	}
}

func SetDefaultLogger(logPath string, logType string) {
	var logger *slog.Logger
	var err error
	if logPath != "" {
		var file *os.File
		file, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panicErr := fmt.Errorf("could not open file %v for logging: %w", logPath, err)
			panic(panicErr)
		}

		logger, err = newLogger(logType, file)
	} else {
		logger, err = newLogger(logType, os.Stdout)
	}

	if err != nil {
		panicErr := fmt.Errorf("could not create logger: %w", err)
		panic(panicErr)
	}
	slog.SetDefault(logger)
}
