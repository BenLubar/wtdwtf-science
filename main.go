package main // import "github.com/BenLubar/wtdwtf-science"

import (
	"context"
	"flag"
	"log"

	"golang.org/x/sync/errgroup"

	"github.com/BenLubar/wtdwtf-science/discourse"
	"github.com/BenLubar/wtdwtf-science/nodebb"
	"github.com/BenLubar/wtdwtf-science/output"
	"github.com/BenLubar/wtdwtf-science/telligent"
	"github.com/pkg/errors"
)

func initDB(ctx context.Context) error {
	wg, ctx := errgroup.WithContext(ctx)

	wg.Go(func() error {
		return errors.Wrap(telligent.Init(ctx), "community server")
	})
	wg.Go(func() error {
		return errors.Wrap(discourse.Init(ctx), "discourse")
	})
	wg.Go(func() error {
		return errors.Wrap(nodebb.Init(ctx), "nodebb")
	})
	wg.Go(func() error {
		return errors.Wrap(output.Init(ctx), "output")
	})
	return wg.Wait()
}

func main() {
	flag.Parse()

	if err := initDB(context.Background()); err != nil {
		log.Fatalf("Database init error: %+v\n", err)
	}
}
