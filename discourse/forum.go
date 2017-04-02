package discourse // import "github.com/BenLubar/wtdwtf-science/discourse"

import (
	"context"
	"database/sql"
	"flag"

	"github.com/BenLubar/wtdwtf-science/forum"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var dataSourceName = flag.String("discourse", "host=wtdwtf-science-discourse user=postgres dbname=discourse sslmode=disable", "Discourse data source name")

func Dial(ctx context.Context) (forum.Forum, error) {
	db, err := sql.Open("postgres", *dataSourceName)
	if err != nil {
		return nil, errors.Wrap(err, "dial discourse")
	}

	if err = db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, errors.Wrap(err, "ping discourse")
	}

	return &Forum{Shared: forum.Shared{ID: "Discourse"}, db: db}, nil
}

type Forum struct {
	forum.Shared
	db        *sql.DB
	telligent forum.Forum
}

func (f *Forum) Close() error {
	return errors.Wrap(f.db.Close(), "close discourse")
}

func (f *Forum) SetPreviousForums(forums []forum.Forum) {
	if len(forums) != 1 {
		panic("discourse: expected 1 previous forum")
	}
	f.telligent = forums[0]
}
