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

	rows, err := f.db.QueryContext(ctx, `select c.id, coalesce(c.parent_category_id, 0), c.name, c.slug, c.description, c.position, ccf.value from categories c left outer join category_custom_fields ccf on ccf.category_id = c.id where ccf.name = 'import_id';`)
	if f.Check(err) {
		return
	}

	defer func() {
		f.Check(rows.Close())
	}()

	for rows.Next() {
		var c Category
		if f.Check(rows.Scan(&c.id, &c.parent, &c.name, &c.slug, &c.description, &c.order, &c.imported)) {
			return
		}
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

func (c *Category) Imported() map[forum.Forum]int64 {
	if c.imported.Valid {
		return map[forum.Forum]int64{
			c.forum.telligent: c.imported.Int64,
		}
	}
	return nil
}
