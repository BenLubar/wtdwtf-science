package main // import "github.com/BenLubar/wtdwtf-science"

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/BenLubar/wtdwtf-science/discourse"
	"github.com/BenLubar/wtdwtf-science/forum"
	"github.com/BenLubar/wtdwtf-science/nodebb"
	"github.com/BenLubar/wtdwtf-science/output"
	"github.com/BenLubar/wtdwtf-science/telligent"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var dialers = [...]forum.DialFunc{
	telligent.Dial,
	discourse.Dial,
	nodebb.Dial,
}

var forums [len(dialers)]forum.Forum

func initForums(ctx context.Context) error {
	var wg errgroup.Group

	for i, d := range dialers {
		i, d := i, d // shadow
		wg.Go(func() error {
			f, err := d(ctx)
			if err != nil {
				return err
			}
			forums[i] = f
			return nil
		})
	}
	wg.Go(func() error {
		return errors.Wrap(output.Init(ctx), "output")
	})
	if err := wg.Wait(); err != nil {
		return err
	}
	for i, f := range forums {
		f.SetPreviousForums(forums[:i])
		if err := output.AddForum(f); err != nil {
			return errors.Wrapf(err, "adding forum %d (%T)", i, f)
		}
	}
	return nil
}

func main() {
	flag.Parse()

	defer func() {
		for _, f := range forums {
			if f != nil {
				if err := f.Close(); err != nil {
					log.Println(err)
				}
			}
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		<-ch
		cancel()
	}()

	if err := initForums(ctx); err != nil {
		log.Printf("Database init error: %+v", err)
		return
	}

	if err := output.Wait(); err != nil {
		log.Printf("Processing error: %+v", err)
		return
	}

	log.Println("Done.")

	// sleep forever
	<-ctx.Done()
}
