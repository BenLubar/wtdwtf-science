package telligent // import "github.com/BenLubar/wtdwtf-science/telligent"

import (
	"context"
	"database/sql"
	"flag"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/pkg/errors"
)

var dataSourceName = flag.String("telligent", "Server=wtdwtf-science-communityserver;Database=TheDailyWtf;User ID=SA;Password=AnInsecure!Passw0rd;", "Community Server data source name")

var db *sql.DB

func Init(ctx context.Context) error {
	var err error
	if db, err = sql.Open("mssql", *dataSourceName); err != nil {
		return errors.Wrap(err, "connect")
	}

	if err = db.PingContext(ctx); err != nil {
		return errors.Wrap(err, "ping")
	}

	log.Println("Connected to Community Server.")

	return nil
}
