package output

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/BenLubar/wtdwtf-science/forum"
	"github.com/pkg/errors"
)

func insertImported(ctx context.Context, stmt *sql.Stmt, tx *sql.Tx, id int64, imported map[forum.Forum]int64) error {
	for importedForum, importedID := range imported {
		importedForumID, err := getForumID(ctx, importedForum)
		if err != nil {
			return errors.Wrap(err, "get imported forum ID")
		}

		_, err = tx.StmtContext(ctx, stmt).ExecContext(ctx, importedForumID, id, importedID)
		if err != nil {
			return errors.Wrapf(err, "insert import from %q", importedForum.Name())
		}
	}
	return nil
}

var insertRawUserStmt, insertImportedUserStmt *sql.Stmt

func insertUser(ctx context.Context, f forum.Forum, u forum.User) error {
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
	err = tx.StmtContext(ctx, insertRawUserStmt).QueryRowContext(ctx, forumID, u.ID(), u.Login(), u.DisplayName(), u.Email(), u.Slug(), u.CreatedAt(), u.LastSeen(), u.Signature(), u.Location(), u.Bio(), u.WebAddress(), u.DateOfBirth()).Scan(&id)
	if err != nil {
		return errors.Wrapf(err, "insert raw user %q from %q", u.Login(), f.Name())
	}

	err = insertImported(ctx, insertImportedUserStmt, tx, id, u.Imported())
	if err != nil {
		return errors.Wrapf(err, "insert imported user %q from %q", u.Login(), f.Name())
	}

	return errors.Wrap(tx.Commit(), "commit transaction")
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
	err = tx.StmtContext(ctx, insertRawCategoryStmt).QueryRowContext(ctx, forumID, c.ID(), c.Parent(), c.Name(), c.Slug(), c.Order(), fmt.Sprintf("%x", c.FgColor()), fmt.Sprintf("%x", c.BgColor())).Scan(&id)
	if err != nil {
		return errors.Wrapf(err, "insert raw category %q from %q", c.Name(), f.Name())
	}

	err = insertImported(ctx, insertImportedCategoryStmt, tx, id, c.Imported())
	if err != nil {
		return errors.Wrapf(err, "insert imported category %q from %q", c.Name(), f.Name())
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
