package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/georgesofianosgr/sizewise/internal/server"
	"github.com/georgesofianosgr/sizewise/pkg/config"
	"github.com/georgesofianosgr/sizewise/pkg/logger"
)

type Flags struct {
	shouldServe bool
	configPath  string
	port        string
	logFile     string
	logType     string
}

func parseFlags() Flags {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panicErr := fmt.Errorf("could not get home directory: %w", err)
		panic(panicErr)
	}
	defaultConfigPath := homeDir + "/" + config.FileName

	var flags Flags
	flag.BoolVar(&flags.shouldServe, "serve", false, "start the server")
	flag.StringVar(&flags.configPath, "config", defaultConfigPath, "config file")
	flag.StringVar(&flags.port, "port", "8080", "port")
	flag.StringVar(&flags.logFile, "log", "", "a path to append logs instead of stdout, ex /var/log/sizewise.log")
	flag.StringVar(&flags.logType, "output", "", "default text, can be 'text' or 'json'")
	flag.Parse()
	return flags
}

func main() {
	flags := parseFlags()
	logger.SetDefaultLogger(flags.logFile, flags.logType)
	slog.Info("Starting server with config file", "config_path", flags.configPath)

	if flags.shouldServe {
		conf, err := config.NewFromFle(flags.configPath)
		if err != nil {
			slog.Error("could not parse config file", "error", err)
			os.Exit(1)
		}

		err = server.Start(conf, flags.port)
		if err != nil {
			slog.Error("could not start server", "error", err)
			os.Exit(1)

		}
	} else {
		flag.PrintDefaults()
	}
}
