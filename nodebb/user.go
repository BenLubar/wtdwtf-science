package nodebb // import "github.com/BenLubar/wtdwtf-science/nodebb"

import (
	"context"
	"strconv"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/BenLubar/wtdwtf-science/forum"
	"github.com/lib/pq"
)

type User struct {
	UID          int64  `bson:"uid"`
	ULogin       string `bson:"username"`
	UDisplayName string `bson:"fullname"`
	UEmail       string `bson:"email"`
	USlug        string `bson:"userslug"`
	UCreatedAt   int64  `bson:"joindate"`
	ULastSeen    int64  `bson:"lastonline"`
	USignature   string `bson:"signature"`
	ULocation    string `bson:"location"`
	UBio         string `bson:"aboutme"`
	UWebAddress  string `bson:"website"`
	UDateOfBirth string `bson:"birthday"`
	discourse    int64
	telligent    int64
	forum        *Forum
}

func (f *Forum) Users(ctx context.Context) <-chan forum.User {
	ch := make(chan forum.User)

	go f.users(ctx, ch)

	return ch
}

func (f *Forum) users(ctx context.Context, ch chan<- forum.User) {
	defer close(ch)

	iter := f.db.DB("0").C("objects").Find(bson.M{
		"_key": "users:joindate",
	}).Iter()

	defer func() {
		f.Check(iter.Close(), "query users")
	}()

	var uid sortedSet

	for iter.Next(&uid) {
		var u *User
		if f.Check(f.db.DB("0").C("objects").Find(bson.M{
			"_key": "user:" + uid.Value,
		}).One(&u), "query user ID %q", uid.Value) {
			return
		}
		var discourse sortedSet
		if err := f.db.DB("0").C("objects").Find(bson.M{
			"_key":  "_imported:_users",
			"score": u.UID,
		}).One(&discourse); err == nil {
			u.discourse, err = strconv.ParseInt(discourse.Value, 10, 64)
			if err != nil {
				f.Check(err, "parse Discourse user ID %d", u.UID)
				u.discourse = 0
			}
		} else if err != mgo.ErrNotFound {
			f.Check(err, "query Discourse user ID %d", u.UID)
			return
		}
		if u.discourse != 0 {
			var telligent sortedSet
			if err := f.db.DB("0").C("objects").Find(bson.M{
				"_key":  "_telligent:_users",
				"score": u.discourse,
			}).One(&telligent); err == nil {
				u.telligent, err = strconv.ParseInt(telligent.Value, 10, 64)
				if err != nil {
					f.Check(err, "parse Community Server user ID %d", u.UID)
					u.telligent = 0
				}
			} else if err != mgo.ErrNotFound {
				f.Check(err, "query Community Server user ID %d", u.UID)
				return
			}
		}
		u.forum = f
		select {
		case ch <- u:
		case <-ctx.Done():
			return
		}
	}
}

func (u *User) ID() int64 {
	return u.UID
}

func (u *User) Login() string {
	return u.ULogin
}

func (u *User) DisplayName() string {
	return u.UDisplayName
}

func (u *User) Email() string {
	return u.UEmail
}

func (u *User) Slug() string {
	return u.USlug
}

func (u *User) CreatedAt() time.Time {
	return makeTime(u.UCreatedAt)
}

func (u *User) LastSeen() pq.NullTime {
	if u.ULastSeen == 0 {
		return pq.NullTime{Valid: false}
	}
	return pq.NullTime{Valid: true, Time: makeTime(u.ULastSeen)}
}

func (u *User) Signature() string {
	return u.USignature
}

func (u *User) Location() string {
	return u.ULocation
}

func (u *User) Bio() string {
	return u.UBio
}

func (u *User) WebAddress() string {
	return u.UWebAddress
}

func (u *User) DateOfBirth() pq.NullTime {
	t, err := time.ParseInLocation("01/02/2006", u.UDateOfBirth, time.UTC)
	if err != nil {
		return pq.NullTime{Valid: false}
	}
	return pq.NullTime{Valid: true, Time: t}
}

func (u *User) Imported() map[forum.Forum]int64 {
	return u.forum.imported(u.discourse, u.telligent)
}
