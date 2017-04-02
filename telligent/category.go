package telligent // import "github.com/BenLubar/wtdwtf-science/telligent"

import (
	"context"

	"github.com/BenLubar/wtdwtf-science/forum"
)

type Category struct {
	id          int64
	parent      int64
	name        string
	description string
	order       int
}

func (f *Forum) Categories(ctx context.Context) <-chan forum.Category {
	ch := make(chan forum.Category)

	go f.categories(ctx, ch)

	return ch
}

func (f *Forum) categories(ctx context.Context, ch chan<- forum.Category) {
	defer close(ch)

	rows, err := f.db.QueryContext(ctx, `select -GroupID, 0, Name, Description, SortOrder from cs_Groups where ApplicationType = 0
union all
select s.SectionID, coalesce(-s.GroupID, 0), s.Name, s.Description, s.SortOrder from cs_Sections s
inner join cs_Groups g on s.GroupID = g.GroupID
where g.ApplicationType = 0`)
	if f.Check(err) {
		return
	}

	defer func() {
		f.Check(rows.Close())
	}()

	for rows.Next() {
		var c Category
		if f.Check(rows.Scan(&c.id, &c.parent, &c.name, &c.description, &c.order)) {
			return
		}
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
	return ""
}

func (c *Category) Description() string {
	return c.description
}

func (c *Category) Order() int {
	return c.order
}

func (c *Category) Imported() map[forum.Forum]int64 {
	return nil
}
