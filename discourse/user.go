package discourse // import "github.com/BenLubar/wtdwtf-science/discourse"

import (
	"context"
	"database/sql"
	"time"

	"github.com/BenLubar/wtdwtf-science/forum"
	"github.com/lib/pq"
)

type User struct {
	id          int64
	login       string
	displayName string
	email       string
	slug        string
	createdAt   time.Time
	lastSeen    pq.NullTime
	location    string
	bio         string
	webAddress  string
	dateOfBirth pq.NullTime
	imported    sql.NullInt64
	forum       *Forum
}

func (f *Forum) Users(ctx context.Context) <-chan forum.User {
	ch := make(chan forum.User)

	go f.users(ctx, ch)

	return ch
}

func (f *Forum) users(ctx context.Context, ch chan<- forum.User) {
	defer close(ch)

	rows, err := f.db.QueryContext(ctx, `select u.id, u.username, coalesce(u.name, u.username), coalesce(u.email, ''), coalesce(u.username_lower, ''), u.created_at, u.last_seen_at, coalesce(p.location, ''), coalesce(p.bio_raw, ''), coalesce(p.website, ''), u.date_of_birth, ucf.value::bigint from users u left outer join user_profiles p on u.id = p.user_id left outer join user_custom_fields ucf on ucf.user_id = u.id and ucf.name = 'import_id';`)
	if f.Check(err, "query users") {
		return
	}

	defer func() {
		f.Check(rows.Close(), "close user query")
	}()

	for rows.Next() {
		var u User
		if f.Check(rows.Scan(&u.id, &u.login, &u.displayName, &u.email, &u.slug, &u.createdAt, &u.lastSeen, &u.location, &u.bio, &u.webAddress, &u.dateOfBirth, &u.imported), "scan user") {
			return
		}
		u.forum = f
		select {
		case ch <- &u:
		case <-ctx.Done():
			return
		}
	}
}

func (u *User) ID() int64 {
	return u.id
}

func (u *User) Login() string {
	return u.login
}

func (u *User) DisplayName() string {
	return u.displayName
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Slug() string {
	return u.slug
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) LastSeen() pq.NullTime {
	return u.lastSeen
}

func (u *User) Signature() string {
	return ""
}

func (u *User) Location() string {
	return u.location
}

func (u *User) Bio() string {
	return u.bio
}

func (u *User) WebAddress() string {
	return u.webAddress
}

func (u *User) DateOfBirth() pq.NullTime {
	return u.dateOfBirth
}

func (u *User) Imported() map[forum.Forum]int64 {
	if u.imported.Valid {
		return map[forum.Forum]int64{
			u.forum.telligent: u.imported.Int64,
		}
	}
	return nil
}
