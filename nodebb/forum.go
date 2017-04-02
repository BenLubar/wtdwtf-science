package nodebb // import "github.com/BenLubar/wtdwtf-science/nodebb"

import (
	"context"
	"flag"

	"github.com/BenLubar/wtdwtf-science/forum"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
)

var dataSourceName = flag.String("nodebb", "wtdwtf-science-nodebb/0", "NodeBB data source name")

func Dial(ctx context.Context) (forum.Forum, error) {
	db, err := mgo.Dial(*dataSourceName)
	if err != nil {
		return nil, errors.Wrap(err, "dial nodebb")
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, errors.Wrap(err, "ping nodebb")
	}

	return &Forum{Shared: forum.Shared{ID: "NodeBB"}, db: db}, nil
}

type Forum struct {
	forum.Shared
	db        *mgo.Session
	telligent forum.Forum
	discourse forum.Forum
}

func (f *Forum) Close() error {
	f.db.Close()
	return nil
}

func (f *Forum) SetPreviousForums(forums []forum.Forum) {
	if len(forums) != 2 {
		panic("nodebb: expected 2 previous forums")
	}
	f.telligent = forums[0]
	f.discourse = forums[1]
}

type sortedSet struct {
	Value string  `json:"value"`
	Score float64 `json:"score"`
}
