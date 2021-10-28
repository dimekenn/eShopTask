package main

import (
	"context"
	"embed"
	"encoding/json"
	"log"
	"os"
	"taskESchop/internal"
	"taskESchop/internal/configs"
)

//go:embed config.json
var fs embed.FS

const configName = "config.json"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//reading json file for configs
	data, readErr := fs.ReadFile(configName)
	if readErr != nil {
		log.Fatal(readErr)
	}

	//creating config entity to deserialize configs.json
	cfg := configs.NewConfig()

	if unmErr := json.Unmarshal(data, &cfg); unmErr != nil {
		log.Fatal(unmErr)
	}

	cfg.DB.Password = os.Getenv("db_password")

	errCh := make(chan error, 1)

	go internal.StartHTTPServer(ctx, errCh, cfg)

	log.Fatalf("%v", <-errCh)
}
