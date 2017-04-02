package nodebb

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/BenLubar/wtdwtf-science/forum"
	"github.com/pkg/errors"
)

type Category struct {
	CID          int64  `bson:"cid"`
	CParent      int64  `bson:"parentCid"`
	CName        string `bson:"name"`
	CSlug        string `bson:"slug"`
	CDescription string `bson:"description"`
	COrder       int    `bson:"order"`
	discourse    int64
	telligent    int64
	forum        *Forum
}

func (f *Forum) Categories(ctx context.Context) <-chan forum.Category {
	ch := make(chan forum.Category)

	go f.categories(ctx, ch)

	return ch
}

func (f *Forum) categories(ctx context.Context, ch chan<- forum.Category) {
	defer close(ch)

	iter := f.db.DB("0").C("objects").Find(bson.M{
		"_key": "categories:cid",
	}).Iter()

	defer func() {
		f.Check(errors.Wrap(iter.Close(), "nodebb categories"))
	}()

	var cid sortedSet

	for iter.Next(&cid) {
		var c *Category
		if f.Check(errors.Wrapf(f.db.DB("0").C("objects").Find(bson.M{
			"_key": "category:" + cid.Value,
		}).One(&c), "nodebb category ID %q", cid.Value)) {
			return
		}
		var discourse sortedSet
		if err := f.db.DB("0").C("objects").Find(bson.M{
			"_key":  "_imported:_categories",
			"score": c.CID,
		}).One(&discourse); err == nil {
			c.discourse, err = strconv.ParseInt(discourse.Value, 10, 64)
			if err != nil {
				f.Check(err)
				c.discourse = 0
			}
		} else if err != mgo.ErrNotFound {
			f.Check(errors.Wrapf(err, "nodebb -> discourse for category ID %d", c.CID))
			return
		}
		if c.discourse != 0 {
			var telligent sortedSet
			if err := f.db.DB("0").C("objects").Find(bson.M{
				"_key":  "_telligent:_categories",
				"score": c.discourse,
			}).One(&telligent); err == nil {
				c.telligent, err = strconv.ParseInt(telligent.Value, 10, 64)
				if err != nil {
					f.Check(err)
					c.telligent = 0
				}
			} else if err != mgo.ErrNotFound {
				f.Check(errors.Wrapf(err, "nodebb -> telligent for category ID %d", c.CID))
				return
			}
		}
		c.forum = f
		c.CSlug = strings.TrimPrefix(c.CSlug, fmt.Sprintf("%d/", c.CID))
		select {
		case ch <- c:
		case <-ctx.Done():
			return
		}
	}
}

func (c *Category) ID() int64 {
	return c.CID
}

func (c *Category) Parent() int64 {
	return c.CParent
}

func (c *Category) Name() string {
	return c.CName
}

func (c *Category) Slug() string {
	return c.CSlug
}

func (c *Category) Description() string {
	return c.CDescription
}

func (c *Category) Order() int {
	return c.COrder
}

func (c *Category) Imported() map[forum.Forum]int64 {
	if c.discourse == 0 && c.telligent == 0 {
		return nil
	}
	m := make(map[forum.Forum]int64)
	if c.discourse != 0 {
		m[c.forum.discourse] = c.discourse
	}
	if c.telligent != 0 {
		m[c.forum.telligent] = c.telligent
	}
	return m
}
