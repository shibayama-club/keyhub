package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/pressly/goose/v3"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flag.String("dir", ".", "directory with migration files")
)

func main() {
	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("goose: failed to parse flags: %v", err)
	}
	args := flags.Args()

	if len(args) < 3 {
		flags.Usage()
		return
	}

	dns, cmd := args[1], args[2]

	db, err := goose.OpenDBWithDriver("postgres", dns)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v", err)
		}
	}()

	arguments := []string{}
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	ctx := context.Background()
	if err := goose.RunContext(ctx, cmd, db, *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", cmd, err)
	}
}
