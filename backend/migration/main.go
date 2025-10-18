package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/pressly/goose/v3"
)

const (
	minRequiredArgs = 3
	dsnArgIndex     = 1
	cmdArgIndex     = 2
	extraArgsOffset = 3
)

var (
	flagSet = flag.NewFlagSet("goose", flag.ExitOnError)
	dirFlag = flagSet.String("dir", ".", "directory with migration files")
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	if err := flagSet.Parse(os.Args[1:]); err != nil {
		return fmt.Errorf("failed to parse flags: %w", err)
	}
	args := flagSet.Args()

	if len(args) < minRequiredArgs {
		flagSet.Usage()
		return fmt.Errorf("insufficient arguments: expected at least %d, got %d", minRequiredArgs, len(args))
	}

	dsn := args[dsnArgIndex]
	command := args[cmdArgIndex]

	db, err := goose.OpenDBWithDriver("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Warning: failed to close database connection: %v", closeErr)
		}
	}()

	var additionalArgs []string
	if len(args) > extraArgsOffset {
		additionalArgs = args[extraArgsOffset:]
	}

	ctx := context.Background()
	if err := goose.RunContext(ctx, command, db, *dirFlag, additionalArgs...); err != nil {
		return fmt.Errorf("goose command %q failed: %w", command, err)
	}

	return nil
}
