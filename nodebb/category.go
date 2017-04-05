package nodebb

import (
	"context"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/BenLubar/wtdwtf-science/forum"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Category struct {
	CID          int64  `bson:"cid"`
	CParent      int64  `bson:"parentCid"`
	CName        string `bson:"name"`
	CSlug        string `bson:"slug"`
	CDescription string `bson:"description"`
	COrder       int    `bson:"order"`
	CFgColor     string `bson:"color"`
	CBgColor     string `bson:"bgColor"`
	CLink        string `bson:"link"`
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
		f.Check(iter.Close(), "query categories")
	}()

	var cid sortedSet

	for iter.Next(&cid) {
		var c *Category
		if f.Check(f.db.DB("0").C("objects").Find(bson.M{
			"_key": "category:" + cid.Value,
		}).One(&c), "query category ID %q", cid.Value) {
			return
		}
		var discourse sortedSet
		if err := f.db.DB("0").C("objects").Find(bson.M{
			"_key":  "_imported:_categories",
			"score": c.CID,
		}).One(&discourse); err == nil {
			c.discourse, err = strconv.ParseInt(discourse.Value, 10, 64)
			if err != nil {
				f.Check(err, "parse Discourse category ID %d", c.CID)
				c.discourse = 0
			}
		} else if err != mgo.ErrNotFound {
			f.Check(err, "query Discourse category ID %d", c.CID)
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
					f.Check(err, "parse Community Server category ID %d", c.CID)
					c.telligent = 0
				}
			} else if err != mgo.ErrNotFound {
				f.Check(err, "query Community Server category ID %d", c.CID)
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

func parseColor(c string, def uint8) [3]uint8 {
	if len(c) == 0 {
		return [3]uint8{def, def, def}
	}
	if c[0] == '#' {
		c = c[1:]
	}
	if len(c) == 3 {
		c = c[:1] + c[:1] + c[1:2] + c[1:2] + c[2:] + c[2:]
	}
	if len(c) != 6 {
		return [3]uint8{def, def, def}
	}
	b, err := hex.DecodeString(c)
	if err != nil || len(b) != 3 {
		return [3]uint8{def, def, def}
	}
	return [3]uint8{b[0], b[1], b[2]}
}

func (c *Category) FgColor() [3]uint8 {
	return parseColor(c.CFgColor, 0xff)
}

func (c *Category) BgColor() [3]uint8 {
	return parseColor(c.CFgColor, 0x00)
}

func (c *Category) Link() string {
	return c.CLink
}

func (c *Category) Imported() map[forum.Forum]int64 {
	return c.forum.imported(c.discourse, c.telligent)
}
