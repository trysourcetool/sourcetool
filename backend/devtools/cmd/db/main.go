package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	"github.com/trysourcetool/sourcetool/backend/internal/config"
	"github.com/trysourcetool/sourcetool/backend/internal/infra/postgres"
	"github.com/trysourcetool/sourcetool/backend/internal/logger"
)

func init() {
	config.Init()
	logger.Init()
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "usage: db [cmd] [args...]\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Commands:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  migrate [dir]: runs all migrations (default dir: migrations)\n")
	}

	flag.Parse()
	if flag.NArg() < 1 {
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
		dir := "migrations"
		if flag.NArg() > 1 {
			dir = flag.Arg(1)
		}
		return migrate(dir)
	default:
		return fmt.Errorf("unsupported arg: %q", cmd)
	}
}

func migrate(dir string) error {
	return postgres.Migrate(dir)
}
