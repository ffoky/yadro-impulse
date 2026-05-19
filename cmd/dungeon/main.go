package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"yadro-impulse/internal/application/coordinator"
	"yadro-impulse/internal/domain"
	"yadro-impulse/internal/infrastructure/config"
	"yadro-impulse/internal/infrastructure/eventsreader"
	"yadro-impulse/internal/infrastructure/output"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	cfgPath := flag.String("config", "config.json", "path to config json")
	eventsPath := flag.String("events", "events", "path to events file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}
	dungeon, err := domain.NewDungeon(cfg)
	if err != nil {
		return fmt.Errorf("create dungeon: %w", err)
	}

	f, err := os.Open(*eventsPath)
	if err != nil {
		return fmt.Errorf("open events file: %w", err)
	}
	defer func() { _ = f.Close() }()

	reader := eventsreader.New(f)
	writer := output.NewWriter(os.Stdout)
	coord := coordinator.New(dungeon, writer)

	for {
		event, readErr := reader.Next()
		if errors.Is(readErr, io.EOF) {
			break
		}
		if readErr != nil {
			return fmt.Errorf("read event file: %w", readErr)
		}
		if handleErr := coord.Handle(event); handleErr != nil {
			return fmt.Errorf("handle event: %w", handleErr)
		}
	}

	coord.CloseDungeon()
	if writeErr := output.WriteReport(os.Stdout, coord.Report()); writeErr != nil {
		return fmt.Errorf("write report: %w", writeErr)
	}
	return nil
}
