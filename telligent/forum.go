package telligent // import "github.com/BenLubar/wtdwtf-science/telligent"

import (
	"context"
	"database/sql"
	"flag"

	"github.com/BenLubar/wtdwtf-science/forum"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/pkg/errors"
)

var dataSourceName = flag.String("telligent", "Server=wtdwtf-science-communityserver;Database=TheDailyWtf;User ID=SA;Password=AnInsecure!Passw0rd;", "Community Server data source name")

func Dial(ctx context.Context) (forum.Forum, error) {
	db, err := sql.Open("mssql", *dataSourceName)
	if err != nil {
		return nil, errors.Wrap(err, "dial telligent")
	}

	if err = db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, errors.Wrap(err, "ping telligent")
	}

	return &Forum{Shared: forum.Shared{ID: "Community Server"}, db: db}, nil
}

type Forum struct {
	forum.Shared
	db *sql.DB
}

func (f *Forum) Close() error {
	return errors.Wrap(f.db.Close(), "close telligent")
}

func (f *Forum) SetPreviousForums(forums []forum.Forum) {
	if len(forums) != 0 {
		panic("telligent: expected 0 previous forums")
	}
}
