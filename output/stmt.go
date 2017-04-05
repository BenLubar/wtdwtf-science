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
		query: `insert into raw_categories (forum_id, raw_id, raw_parent_id, name, slug, position, fg_color, bg_color) values ($1::bigint, $2::bigint, nullif($3::bigint, 0), $4::text, $5::text, $6::int, $7::char(6), $8::char(6)) returning id;`,
	},
	{
		stmt:  &insertImportedCategoryStmt,
		query: `insert into imported_categories (imported_forum_id, category_id, imported_category_id) values ($1::bigint, $2::bigint, $3::bigint);`,
	},
	{
		stmt:  &insertRawUserStmt,
		query: `insert into raw_users (forum_id, raw_id, login, display_name, email, slug, created_at, last_seen, signature, location, bio, web_address, date_of_birth) values ($1::bigint, $2::bigint, $3::text, $4::text, $5::text, $6::text, $7::timestamptz, $8::timestamptz, $9::text, $10::text, $11::text, $12::text, $13::date) returning id;`,
	},
	{
		stmt:  &insertImportedUserStmt,
		query: `insert into imported_users (imported_forum_id, user_id, imported_user_id) values ($1::bigint, $2::bigint, $3::bigint);`,
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
