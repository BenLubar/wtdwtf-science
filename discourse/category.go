package discourse // import "github.com/BenLubar/wtdwtf-science/discourse"

import (
	"context"
	"database/sql"

	"github.com/BenLubar/wtdwtf-science/forum"
)

type Category struct {
	id          int64
	parent      int64
	name        string
	slug        string
	description string
	order       int
	fgColor     [3]uint8
	bgColor     [3]uint8
	link        string
	imported    sql.NullInt64
	forum       *Forum
}

func (f *Forum) Categories(ctx context.Context) <-chan forum.Category {
	ch := make(chan forum.Category)

	go f.categories(ctx, ch)

	return ch
}

func (f *Forum) categories(ctx context.Context, ch chan<- forum.Category) {
	defer close(ch)

	rows, err := f.db.QueryContext(ctx, `select c.id, coalesce(c.parent_category_id, 0), c.name, c.slug, c.description, c.position, decode(coalesce(c.text_color, 'ffffff'), 'hex'), decode(coalesce(c.color, '000000'), 'hex'), ccf.value from categories c left outer join category_custom_fields ccf on ccf.category_id = c.id and ccf.name = 'import_id';`)
	if f.Check(err, "query categories") {
		return
	}

	defer func() {
		f.Check(rows.Close(), "close category query")
	}()

	for rows.Next() {
		var c Category
		var fgColor, bgColor []byte
		if f.Check(rows.Scan(&c.id, &c.parent, &c.name, &c.slug, &c.description, &c.order, &fgColor, &bgColor, &c.imported), "scan category") {
			return
		}
		copy(c.fgColor[:], fgColor)
		copy(c.bgColor[:], bgColor)
		c.forum = f
		select {
		case ch <- &c:
		case <-ctx.Done():
			return
		}
	}
}

func (c *Category) ID() int64 {
	return c.id
}

func (c *Category) Parent() int64 {
	return c.parent
}

func (c *Category) Name() string {
	return c.name
}

func (c *Category) Slug() string {
	return c.slug
}

func (c *Category) Description() string {
	return c.description
}

func (c *Category) Order() int {
	return c.order
}

func (c *Category) FgColor() [3]uint8 {
	return c.fgColor
}

func (c *Category) BgColor() [3]uint8 {
	return c.bgColor
}

func (c *Category) Link() string {
	return c.link
}

func (c *Category) Imported() map[forum.Forum]int64 {
	if c.imported.Valid {
		return map[forum.Forum]int64{
			c.forum.telligent: c.imported.Int64,
		}
	}
	return nil
}
