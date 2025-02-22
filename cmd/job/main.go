package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/change-feed/checker/internal/content"
	"github.com/change-feed/checker/pkg/storage"
)

// Configure the logger
func init() {
	level := func() slog.Level {
		switch strings.ToUpper(os.Getenv("LOG_LEVEL")) {
		case "ERROR":
			return slog.LevelError
		case "WARN":
			return slog.LevelWarn
		case "INFO":
			return slog.LevelInfo
		case "DEBUG":
			return slog.LevelDebug
		default:
			return slog.LevelInfo
		}
	}()
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})))
}

func main() {
	backend, err := storage.NewDaprStorage()
	if err != nil {
		slog.Error("Error initializing storage backend", "error", err)
		os.Exit(1)
	}

	if err := run(backend); err != nil {
		slog.Error("Unable to process!", "error", err)
		os.Exit(1)
	}
}

func run(backend storage.StorageBackend) error {
	url := os.Getenv("CHECK_URL")
	if url == "" {
		return fmt.Errorf("CHECK_URL environment variable is required")
	}

	checker, err := content.NewContentChecker(backend, url)
	if err != nil {
		log.Fatal(err)
	}

	_, err = checker.Check()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
