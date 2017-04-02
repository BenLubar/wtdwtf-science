package output

import (
	"context"
	"database/sql"

	"github.com/BenLubar/wtdwtf-science/forum"
	"github.com/pkg/errors"
)

func insertUser(ctx context.Context, f forum.Forum, u forum.User) error {
	return nil
}

func insertGroup(ctx context.Context, f forum.Forum, g forum.Group) error {
	return nil
}

var insertRawCategoryStmt, insertImportedCategoryStmt *sql.Stmt

func insertCategory(ctx context.Context, f forum.Forum, c forum.Category) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "begin transaction")
	}
	defer tx.Rollback()

	forumID, err := getForumID(ctx, f)
	if err != nil {
		return errors.Wrap(err, "get forum ID")
	}

	var id int64
	err = tx.StmtContext(ctx, insertRawCategoryStmt).QueryRowContext(ctx, forumID, c.ID(), c.Parent(), c.Name(), c.Slug(), c.Order()).Scan(&id)
	if err != nil {
		return errors.Wrapf(err, "insert raw category %q from %q", c.Name(), f.Name())
	}

	for importedForum, importedID := range c.Imported() {
		importedForumID, err := getForumID(ctx, importedForum)
		if err != nil {
			return errors.Wrap(err, "get imported forum ID")
		}

		_, err = tx.StmtContext(ctx, insertImportedCategoryStmt).ExecContext(ctx, importedForumID, id, importedID)
		if err != nil {
			return errors.Wrapf(err, "insert category %q imported %q -> %q", c.Name(), importedForum.Name(), f.Name())
		}
	}

	return errors.Wrap(tx.Commit(), "commit transaction")
}

func insertTopic(ctx context.Context, f forum.Forum, t forum.Topic) error {
	return nil
}

func insertPost(ctx context.Context, f forum.Forum, p forum.Post) error {
	return nil
}

func insertVote(ctx context.Context, f forum.Forum, v forum.Vote) error {
	return nil
}
