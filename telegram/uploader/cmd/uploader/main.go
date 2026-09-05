package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"uploader/internal/config"
	"uploader/internal/uploader"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Fatal: failed to load config: %v\n", err)
		os.Exit(1)
	}

	store, err := uploader.NewStore(cfg.DatabaseURL)
	if err != nil {
		fmt.Printf("Fatal: failed to initialize store: %v\n", err)
		os.Exit(1)
	}
	defer store.Close()

	u := uploader.New(cfg, store)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.APITimeout)*time.Second)
	defer cancel()

	if err := u.Run(ctx); err != nil {
		fmt.Printf("Fatal: %v\n", err)
		os.Exit(1)
	}
}
