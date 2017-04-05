package output // import "github.com/BenLubar/wtdwtf-science/output"

import (
	"context"

	"github.com/lib/pq"
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

	if _, err := db.ExecContext(ctx, `drop schema public cascade;`); err != nil {
		// a missing schema is ok
		if perr, ok := err.(*pq.Error); !ok || perr.Code != "3F000" {
			panic(initError{errors.Wrap(err, "init schema\ndrop schema public cascade;")})
		}
	}
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
	fg_color char(6) not null,
	bg_color char(6) not null,

	unique (forum_id, raw_id),
	check(fg_color ~ '^[0-9a-f]{6}$'),
	check(bg_color ~ '^[0-9a-f]{6}$')
);`)
	do(`create table imported_categories (
	imported_forum_id bigint not null references forums (id),
	category_id bigint not null references raw_categories (id),
	imported_category_id bigint not null,

	unique (imported_forum_id, category_id)
);`)
	do(`create table raw_users (
	id bigserial primary key,
	forum_id bigint not null references forums (id),
	raw_id bigint not null,
	login text not null,
	display_name text not null,
	email text not null,
	slug text not null,
	created_at timestamptz not null,
	last_seen timestamptz,
	signature text not null,
	location text not null,
	bio text not null,
	web_address text not null,
	date_of_birth date,

	unique (forum_id, raw_id)
);`)
	do(`create table imported_users (
	imported_forum_id bigint not null references forums (id),
	user_id bigint not null references raw_users (id),
	imported_user_id bigint not null,

	unique (imported_forum_id, user_id)
);`)

	return nil
}
