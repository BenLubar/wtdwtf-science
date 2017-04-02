package nodebb // import "github.com/BenLubar/wtdwtf-science/nodebb"

import (
	"context"
	"flag"
	"log"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
)

var dataSourceName = flag.String("nodebb", "wtdwtf-science-nodebb/0", "NodeBB data source name")

var db *mgo.Session

func Init(ctx context.Context) error {
	var err error
	if db, err = mgo.Dial(*dataSourceName); err != nil {
		return errors.Wrap(err, "connect")
	}

	if err = db.Ping(); err != nil {
		return errors.Wrap(err, "ping")
	}

	log.Println("Connected to NodeBB.")

	return nil
}
