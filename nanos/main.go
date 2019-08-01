package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/olivere/env"
	"github.com/pkg/errors"
)

func main() {
	if err := runMain(); err != nil {
		log.Fatal(err)
	}
}

func runMain() error {
	var (
		urn = flag.String("urn", env.String("", "URN"), "MySQL database connection string; see https://github.com/go-sql-driver/mysql#usage for details")
	)
	flag.Parse()

	if flag.NArg() == 0 {
		return errors.New("please specify a command like list or insert")
	}

	if *urn == "" {
		return errors.New("missing database connection string; use -urn to pass it")
	}

	switch cmd := strings.ToLower(flag.Arg(0)); cmd {
	default:
		return errors.Errorf("unknown command: %s", cmd)
	case "list":
		r, err := NewMySQLRepository(*urn)
		if err != nil {
			return err
		}
		return list(context.Background(), r)
	}
}

func list(ctx context.Context, r Repository) error {
	people, err := r.SelectAll(ctx)
	if err != nil {
		return err
	}
	for _, p := range people {
		fmt.Fprintf(os.Stdout, "%5d  %-20s  %-30s\n",
			p.ID,
			p.Name,
			p.Created.Format(time.RFC3339Nano),
		)
	}
	return nil
}
