package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/trysourcetool/sourcetool/backend/config"
	"github.com/trysourcetool/sourcetool/backend/logger"
	"github.com/trysourcetool/sourcetool/backend/postgres"
)

func init() {
	config.Init()
	logger.Init()
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: db [cmd]\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  migrate: runs all migrations \n")
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	ctx := context.Background()
	if err := run(ctx, flag.Args()[0]); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context, cmd string) error {
	switch cmd {
	case "migrate":
		return migrate()
	default:
		return fmt.Errorf("unsupported arg: %q", cmd)
	}
}

func migrate() error {
	return postgres.Migrate()
}
