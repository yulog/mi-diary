package infra

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/model"
	mi "github.com/yulog/miutil"
)

// func (infra *Infra) InsertFromFile(ctx context.Context, profile string) {
// 	// JSON読み込み
// 	f, _ := os.ReadFile("users_reactions.json")
// 	var r mi.Reactions
// 	json.Unmarshal(f, &r)

// 	infra.Insert(ctx, profile, &r)
// }

func (infra *Infra) RunInTx(ctx context.Context, profile string, fn func(ctx context.Context, tx bun.Tx) error) {
	err := infra.DB(profile).RunInTx(ctx, &sql.TxOptions{}, fn)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}

func (infra *Infra) InsertHashTag(ctx context.Context, db bun.IDB, hashtag *model.HashTag) error {
	_, err := db.NewInsert().Model(hashtag).On("CONFLICT DO UPDATE").Exec(ctx)
	return err
}

func (infra *Infra) InsertNotes(ctx context.Context, db bun.IDB, notes *[]model.Note) (int64, error) {
	result, err := db.NewInsert().Model(notes).Ignore().Exec(ctx)
	rows, _ := result.RowsAffected()
	return rows, err
}

func (infra *Infra) InsertReactions(ctx context.Context, db bun.IDB, reactions *[]model.ReactionEmoji) error {
	_, err := db.NewInsert().Model(reactions).Ignore().Exec(ctx)
	return err
}

func (infra *Infra) InsertNoteToTags(ctx context.Context, db bun.IDB, noteToTags *[]model.NoteToTag) error {
	_, err := db.NewInsert().Model(noteToTags).Ignore().Exec(ctx)
	return err
}

func (infra *Infra) InsertNoteToFiles(ctx context.Context, db bun.IDB, noteToFiles *[]model.NoteToFile) error {
	_, err := db.NewInsert().Model(noteToFiles).Ignore().Exec(ctx)
	return err
}

func (infra *Infra) Count(ctx context.Context, db bun.IDB) error {
	// リアクションのカウント
	err := countReaction(ctx, db)
	if err != nil {
		return err
	}

	// タグのカウント
	err = countHashTag(ctx, db)
	if err != nil {
		return err
	}

	// ユーザーのカウント
	err = countUser(ctx, db)
	if err != nil {
		return err
	}

	// 月別のカウント
	err = countMonthly(ctx, db)
	if err != nil {
		return err
	}

	// 日別のカウント
	err = countDaily(ctx, db)
	if err != nil {
		return err
	}
	return err
}

// リアクションのカウント
func countReaction(ctx context.Context, db bun.IDB) error {
	var reactions []model.ReactionEmoji
	err := db.NewSelect().
		Model((*model.Note)(nil)).
		Relation("Reaction", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.ColumnExpr("reaction.id as id")
		}).
		ColumnExpr("reaction_emoji_name as name").
		ColumnExpr("count(*) as count").
		Group("reaction_emoji_name").
		Scan(ctx, &reactions)
	if err != nil {
		return err
	}

	_, err = db.NewUpdate().
		Model(&reactions).
		OmitZero().
		Column("count").
		Bulk().
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// タグのカウント
func countHashTag(ctx context.Context, db bun.IDB) error {
	var hashtags []model.HashTag
	err := db.NewSelect().
		Model((*model.NoteToTag)(nil)).
		Relation("HashTag", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("text")
		}).
		ColumnExpr("hash_tag_id as id").
		ColumnExpr("count(*) as count").
		Group("hash_tag_id").
		Scan(ctx, &hashtags)
	if err != nil {
		return err
	}

	if len(hashtags) > 0 {
		_, err = db.NewUpdate().
			Model(&hashtags).
			OmitZero().
			Column("count").
			Bulk().
			Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// ユーザーのカウント
func countUser(ctx context.Context, db bun.IDB) error {
	var users []model.User
	err := db.NewSelect().
		Model((*model.Note)(nil)).
		Relation("User", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("id", "name")
		}).
		ColumnExpr("count(*) as count").
		Group("user_id").
		Scan(ctx, &users)
	if err != nil {
		return err
	}

	_, err = db.NewUpdate().
		Model(&users).
		OmitZero().
		Column("count").
		Bulk().
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// 月別のカウント
func countMonthly(ctx context.Context, db bun.IDB) error {
	var months []model.Month
	err := db.NewSelect().
		Model((*model.Note)(nil)).
		ColumnExpr("strftime('%Y-%m', created_at, 'localtime') as ym").
		ColumnExpr("count(*) as count").
		Group("ym").
		Having("ym is not null").
		Scan(ctx, &months)
	if err != nil {
		return err
	}

	_, err = db.NewInsert().
		Model(&months).
		On("CONFLICT DO UPDATE").
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// 日別のカウント
func countDaily(ctx context.Context, db bun.IDB) error {
	var days []model.Day
	err := db.NewSelect().
		Model((*model.Note)(nil)).
		ColumnExpr("strftime('%Y-%m-%d', created_at, 'localtime') as ymd").
		ColumnExpr("strftime('%Y-%m', created_at, 'localtime') as ym").
		ColumnExpr("count(*) as count").
		Group("ymd").
		Having("ymd is not null").
		Scan(ctx, &days)
	if err != nil {
		return err
	}

	_, err = db.NewInsert().
		Model(&days).
		On("CONFLICT DO UPDATE").
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (infra *Infra) UpdateEmoji(ctx context.Context, profile string, id int64, e *mi.Emoji) {
	// TODO: emoji画像をローカルに保存する

	r := model.ReactionEmoji{
		ID:    id,
		Image: e.URL,
	}
	var s []model.ReactionEmoji
	s = append(s, r)
	_, err := infra.DB(profile).NewUpdate().
		Model(&s).
		OmitZero().
		Column("image").
		Bulk().
		Exec(ctx)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}

func (infra *Infra) UpdateColor(ctx context.Context, profile, id, c1, c2 string) {
	r := model.File{
		ID:            id,
		DominantColor: c1,
		GroupColor:    c2,
	}
	var s []model.File
	s = append(s, r)
	_, err := infra.DB(profile).NewUpdate().
		Model(&s).
		OmitZero().
		Column("dominant_color").
		Column("group_color").
		Bulk().
		Exec(ctx)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}
