package output // import "github.com/BenLubar/wtdwtf-science/output"

import (
	"context"
	"database/sql"
	"sync"

	"github.com/BenLubar/wtdwtf-science/forum"
	"github.com/pkg/errors"
)

var (
	initData     = make(map[forum.Forum]*forumInitData)
	initDataLock sync.Mutex
)

type forumInitData struct {
	id   int64
	wait chan struct{}
}

var addForumStmt *sql.Stmt

func addForum(ctx context.Context, f forum.Forum) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrapf(err, "begin transaction for add forum %q", f.Name())
	}
	defer tx.Rollback()

	var data forumInitData

	err = tx.StmtContext(ctx, addForumStmt).QueryRowContext(ctx, f.Name()).Scan(&data.id)
	if err != nil {
		return errors.Wrapf(err, "add forum %q", f.Name())
	}

	initDataLock.Lock()
	if d, ok := initData[f]; ok {
		data.wait = d.wait
		*d = data
		close(data.wait)
	} else {
		data.wait = make(chan struct{})
		initData[f] = &data
		close(data.wait)
	}
	initDataLock.Unlock()

	return errors.Wrapf(tx.Commit(), "commit add forum %q", f.Name())
}

func getForumID(ctx context.Context, f forum.Forum) (int64, error) {
	initDataLock.Lock()
	d, ok := initData[f]
	if !ok {
		d = &forumInitData{wait: make(chan struct{})}
		initData[f] = d
	}
	initDataLock.Unlock()

	select {
	case <-d.wait:
		return d.id, nil
	case <-ctx.Done():
		return 0, ctx.Err()
	}
}

func process(ctx context.Context, f forum.Forum) error {
	for u := range f.Users(ctx) {
		if err := insertUser(ctx, f, u); err != nil {
			return errors.Wrapf(err, "add user %q from %q", u.Name(), f.Name())
		}
	}
	for c := range f.Categories(ctx) {
		if err := insertCategory(ctx, f, c); err != nil {
			return errors.Wrapf(err, "add category %q from %q", c.Name(), f.Name())
		}
	}
	for t := range f.Topics(ctx) {
		if err := insertTopic(ctx, f, t); err != nil {
			return errors.Wrapf(err, "add topic %d from %q", t.ID(), f.Name())
		}
	}
	for p := range f.Posts(ctx) {
		if err := insertPost(ctx, f, p); err != nil {
			return errors.Wrapf(err, "add post %d from %q", p.ID(), f.Name())
		}
	}
	for v := range f.Votes(ctx) {
		if err := insertVote(ctx, f, v); err != nil {
			return errors.Wrapf(err, "add vote %d -> %d from %q", v.UserID(), v.PostID(), f.Name())
		}
	}
	return f.Err()
}
