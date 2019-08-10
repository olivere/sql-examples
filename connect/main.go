package main

// Connect to MySQL with gocloud.dev.
//
// Example:
//
//    ./connect -urn=mysql://user:password@localhost/testdb
//    ./connect -urn=mysql://user:password@tcp(127.0.0.1:3306)/testdb
//    ./connect -urn=gcpmysql://user:password@example-project/region/my-instance01/testdb
//    ./connect -urn=awsmysql://myrole:swordfish@example01.xyzzy.us-west-1.rds.amazonaws.com/testdb
//    ./connect -urn=azuremysql://user:password@example00.mysql.database.azure.com/testdb
import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/olivere/env"
	"github.com/pkg/errors"
	"gocloud.dev/mysql"
	_ "gocloud.dev/mysql/awsmysql"
	_ "gocloud.dev/mysql/azuremysql"
	_ "gocloud.dev/mysql/gcpmysql"
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

	if *urn == "" {
		return errors.New("missing -urn parameter")
	}

	// Connect
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db, err := mysql.Open(ctx, *urn)
	if err != nil {
		return errors.Wrap(err, "unable to connect")
	}
	// Ping
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to ping")
	}
	defer db.Close()

	fmt.Println("OK")

	return nil
}
