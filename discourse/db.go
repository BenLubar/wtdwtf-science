package discourse // import "github.com/BenLubar/wtdwtf-science/discourse"

import (
	"context"
	"database/sql"
	"flag"
	"log"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var dataSourceName = flag.String("discourse", "host=wtdwtf-science-discourse user=postgres dbname=postgres sslmode=disable", "Discourse data source name")

var db *sql.DB

func Init(ctx context.Context) error {
	var err error
	if db, err = sql.Open("postgres", *dataSourceName); err != nil {
		return errors.Wrap(err, "connect")
	}

	if err = db.PingContext(ctx); err != nil {
		return errors.Wrap(err, "ping")
	}

	log.Println("Connected to Discourse.")

	return nil
}
