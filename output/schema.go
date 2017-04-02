package output // import "github.com/BenLubar/wtdwtf-science/output"

import (
	"context"

	"github.com/pkg/errors"
)

func initSchema(ctx context.Context) (err error) {
	type initError struct{ err error }
	defer func() {
		if r := recover(); r != nil {
			err = r.(initError).err
		}
	}()
	do := func(query string, args ...interface{}) {
		if _, err := db.ExecContext(ctx, query, args...); err != nil {
			panic(initError{errors.Wrapf(err, "init schema\n%s", query)})
		}
	}

	do(`drop schema public cascade;`)
	do(`create schema public;`)
	do(`grant all on schema public to postgres;`)
	do(`grant all on schema public to public;`)

	do(`create table forums (
	id bigserial primary key,
	name text not null unique
);`)
	do(`create table raw_categories (
	id bigserial primary key,
	forum_id bigint not null references forums (id),
	raw_id bigint not null,
	raw_parent_id bigint,
	name text,
	slug text,
	position int not null,

	unique (forum_id, raw_id)
);`)
	do(`create table imported_categories (
	imported_forum_id bigint not null references forums (id),
	category_id bigint not null references raw_categories (id),
	imported_category_id bigint not null,

	unique (imported_forum_id, category_id)
);`)

	return nil
}
