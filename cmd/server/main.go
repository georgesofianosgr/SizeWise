package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/georgesofianosgr/sizewise/internal/server"
	"github.com/georgesofianosgr/sizewise/pkg/config"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "development" {
		fmt.Println("Running in development mode")
		err := godotenv.Load()
		if err != nil {
			fmt.Printf("Error loading .env file %s\n", err)
		}
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panicErr := fmt.Errorf("could not get home directory: %w", err)
		panic(panicErr)
	}

	var shouldServe bool
	var configPath string
	var port string
	defaultConfigPath := homeDir + "/" + config.FileName
	flag.BoolVar(&shouldServe, "serve", false, "start the server")
	flag.StringVar(&configPath, "config", defaultConfigPath, "config file")
	flag.StringVar(&port, "port", "8080", "port")
	flag.Parse()

	if shouldServe {
		conf, err := config.NewFromFle(configPath)
		if err != nil {
			panicErr := fmt.Errorf("could not parse config file: %w", err)
			panic(panicErr)
		}
		err = server.Start(conf, port)
		if err != nil {
			panicErr := fmt.Errorf("could not start server: %w", err)
			panic(panicErr)
		}
	} else {
		flag.PrintDefaults()
	}
}
