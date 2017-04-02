package output // import "github.com/BenLubar/wtdwtf-science/output"

import (
	"context"
	"database/sql"
	"flag"
	"log"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var dataSourceName = flag.String("output", "host=wtdwtf-science-output user=postgres dbname=postgres sslmode=disable", "output data source name")

var db *sql.DB

func Init(ctx context.Context) error {
	var err error
	if db, err = sql.Open("postgres", *dataSourceName); err != nil {
		return errors.Wrap(err, "connect")
	}

	if err = db.PingContext(ctx); err != nil {
		return errors.Wrap(err, "ping")
	}

	log.Println("Connected to output database.")

	return nil
}
