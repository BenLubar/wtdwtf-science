package output // import "github.com/BenLubar/wtdwtf-science/output"

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

var statements = [...]struct {
	stmt  **sql.Stmt
	query string
}{
	{
		stmt:  &addForumStmt,
		query: `insert into forums (name) values ($1::text) returning id;`,
	},
	{
		stmt:  &insertRawCategoryStmt,
		query: `insert into raw_categories (forum_id, raw_id, raw_parent_id, name, slug, position) values ($1::bigint, $2::bigint, nullif($3::bigint, 0), $4::text, $5::text, $6::int) returning id;`,
	},
	{
		stmt:  &insertImportedCategoryStmt,
		query: `insert into imported_categories (imported_forum_id, category_id, imported_category_id) values ($1::bigint, $2::bigint, $3::bigint);`,
	},
}

func prepareStatements(ctx context.Context) error {
	var err error
	for _, s := range statements {
		*s.stmt, err = db.PrepareContext(ctx, s.query)
		if err != nil {
			return errors.Wrapf(err, "preparing statement: %q", s.query)
		}
	}

	return nil
}
