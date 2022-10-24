package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rgallagher27/porter/internal/services/port"
	"github.com/rgallagher27/porter/internal/store"
	"github.com/rgallagher27/porter/pkg/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(sigChan)
		cancel()
	}()

	go func() {
		select {
		case <-sigChan:
			cancel()
		case <-ctx.Done():
		}
	}()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	if err := Run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}

// Run is the main process of the application.
func Run(ctx context.Context, cfg config.Config) error {
	log.Println("Running with config", cfg)

	str, err := store.New(store.Config{
		Address:  cfg.RedisAddr,
		Password: "",
		DB:       0,
	})
	if err != nil {
		return fmt.Errorf("new store: %w", err)
	}

	rdr, err := openFile(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}

	portService := port.New(str, cfg.IgnoreErrors)

	err = portService.Run(ctx, rdr)
	if err != nil {
		return fmt.Errorf("port service run: %w", err)
	}

	log.Println("Run completed")

	return nil
}

// openFile is a simple helper for opening a local file, returning an io.ReadCloser
// In a real world application this would be a separate package but for times sake a simple version
// was used here
func openFile(fileName string) (io.ReadCloser, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return f, nil
}
