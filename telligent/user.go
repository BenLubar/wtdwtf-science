package telligent // import "github.com/BenLubar/wtdwtf-science/telligent"

import (
	"context"
	"time"

	"github.com/BenLubar/wtdwtf-science/forum"
	"github.com/lib/pq"
)

type User struct {
	id         int64
	login      string
	email      string
	createdAt  time.Time
	lastSeen   pq.NullTime
	signature  string
	location   string
	bio        string
	webAddress string
}

func (f *Forum) Users(ctx context.Context) <-chan forum.User {
	ch := make(chan forum.User)

	go f.users(ctx, ch)

	return ch
}

func (f *Forum) users(ctx context.Context, ch chan<- forum.User) {
	defer close(ch)

	rows, err := f.db.QueryContext(ctx, `select u.UserID, u.UserName, u.Email, u.CreateDate, u.LastActivity, p.PropertyNames, p.PropertyValuesString, cast(p.PropertyValuesBinary as varbinary(max)) from cs_Users u
left outer join aspnet_Profile p on p.UserId = u.MembershipID`)
	if f.Check(err, "query users") {
		return
	}

	defer func() {
		f.Check(rows.Close(), "close user query")
	}()

	for rows.Next() {
		var u User
		var p properties
		if f.Check(rows.Scan(&u.id, &u.login, &u.email, &u.createdAt, &u.lastSeen, &p.keys, &p.strings, &p.bytes), "scan user") {
			return
		}
		u.signature = p.StringOrDefault("signature", "")
		u.location = p.StringOrDefault("location", "")
		u.bio = p.StringOrDefault("bio", "")
		u.webAddress = p.StringOrDefault("webAddress", "")
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
	return u.login
}

func (u *User) Email() string {
	return u.email
}

func (u *User) Slug() string {
	return u.login
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func (u *User) LastSeen() pq.NullTime {
	return u.lastSeen
}

func (u *User) Signature() string {
	return u.signature
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
	return pq.NullTime{Valid: false}
}

func (u *User) Imported() map[forum.Forum]int64 {
	return nil
}
