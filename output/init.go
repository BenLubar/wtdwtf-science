package output // import "github.com/BenLubar/wtdwtf-science/output"

import (
	"context"
	"database/sql"
	"flag"
	"log"

	"golang.org/x/sync/errgroup"

	"github.com/BenLubar/wtdwtf-science/forum"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var dataSourceName = flag.String("output", "host=wtdwtf-science-output user=postgres dbname=postgres sslmode=disable", "output data source name")

var db *sql.DB
var wg *errgroup.Group
var wgCtx context.Context

func Init(ctx context.Context) error {
	wg, wgCtx = errgroup.WithContext(ctx)

	var err error
	if db, err = sql.Open("postgres", *dataSourceName); err != nil {
		return errors.Wrap(err, "dial output")
	}

	if err = db.PingContext(ctx); err != nil {
		_ = db.Close()
		return errors.Wrap(err, "ping output")
	}

	log.Println("Connected to output database.")

	if err = initSchema(ctx); err != nil {
		_ = db.Close()
		return errors.Wrap(err, "init output schema")
	}

	if err = prepareStatements(ctx); err != nil {
		_ = db.Close()
		return errors.Wrap(err, "prepare output statements")
	}

	return nil
}

func AddForum(f forum.Forum) error {
	if err := addForum(wgCtx, f); err != nil {
		return errors.Wrapf(err, "adding %q", f.Name())
	}

	log.Printf("Added forum: %q", f.Name())

	wg.Go(func() error {
		return process(wgCtx, f)
	})

	return nil
}

func Wait() error {
	if err := wg.Wait(); err != nil {
		return err
	}
	return db.Close()
}
