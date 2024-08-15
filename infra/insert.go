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

// func test(r *mi.Reactions, rowsp *int64) func(ctx context.Context, tx bun.Tx) error {
// 	return func(ctx context.Context, tx bun.Tx) error {
// 		// JSONの中身をモデルへ移す
// 		var (
// 			users       []model.User
// 			notes       []model.Note
// 			reactions   []model.ReactionEmoji
// 			noteToTags  []model.NoteToTag
// 			files       []model.File
// 			noteToFiles []model.NoteToFile
// 		)

// 		for _, v := range *r {
// 			var dn string
// 			if v.Note.User.Name == nil {
// 				dn = v.Note.User.Username
// 			} else {
// 				dn = v.Note.User.Name.(string)
// 			}
// 			u := model.User{
// 				ID:          v.Note.User.ID,
// 				Name:        v.Note.User.Username,
// 				DisplayName: dn,
// 				AvatarURL:   v.Note.User.AvatarURL,
// 			}
// 			users = append(users, u)

// 			for _, tv := range v.Note.Tags {
// 				ht := model.HashTag{Text: tv}
// 				err := insertHashTag(ctx, tx, &ht)
// 				if err != nil {
// 					return err
// 				}
// 				slog.Info("HashTag ID", slog.Int64("ID", ht.ID))
// 				// pp.Println(ht.ID)
// 				// id is scanned automatically
// 				noteToTags = append(noteToTags, model.NoteToTag{NoteID: v.Note.ID, HashTagID: ht.ID})
// 			}

// 			for _, fv := range v.Note.Files {
// 				f := model.File{
// 					ID:           fv.ID,
// 					Name:         fv.Name,
// 					URL:          fv.URL,
// 					ThumbnailURL: fv.ThumbnailURL,
// 					Type:         fv.Type,
// 					CreatedAt:    fv.CreatedAt,
// 				}
// 				files = append(files, f)
// 				noteToFiles = append(noteToFiles, model.NoteToFile{NoteID: v.Note.ID, FileID: f.ID})
// 			}

// 			reactionName := strings.TrimSuffix(strings.TrimPrefix(v.Note.MyReaction, ":"), "@.:")
// 			n := model.Note{
// 				ID:                v.Note.ID,
// 				ReactionID:        v.ID,
// 				UserID:            v.Note.User.ID,
// 				ReactionEmojiName: reactionName,
// 				Text:              v.Note.Text,
// 				CreatedAt:         v.Note.CreatedAt, // SQLite は日時をUTCで保持する
// 			}
// 			notes = append(notes, n)

// 			r := model.ReactionEmoji{
// 				Name: reactionName,
// 			}
// 			reactions = append(reactions, r)
// 		}

// 		// 重複していたらアップデート
// 		err := insertUsers(ctx, tx, &users)
// 		if err != nil {
// 			return err
// 		}

// 		// 重複していたら登録しない(エラーにしない)
// 		rows, err := insertNotes(ctx, tx, &notes)
// 		if err != nil {
// 			return err
// 		}
// 		slog.Info("Notes inserted", slog.Int64("count", rows))
// 		*rowsp = rows
// 		slog.Info("Notes inserted(pointer)", slog.Int64("count", *rowsp))

// 		err = insertReactions(ctx, tx, &reactions)
// 		if err != nil {
// 			return err
// 		}

// 		// 0件の場合がある
// 		if len(noteToTags) > 0 {
// 			err = insertNoteToTags(ctx, tx, &noteToTags)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		if len(files) > 0 {
// 			err = insertFiles(ctx, tx, &files)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		if len(noteToFiles) > 0 {
// 			err = insertNoteToFiles(ctx, tx, &noteToFiles)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		err = count(ctx, tx)
// 		// TODO: あってもなくても変わらない vs 統一感
// 		if err != nil {
// 			return err
// 		}
// 		return err
// 	}
// }

// func (infra *Infra) Insert(ctx context.Context, profile string, r *mi.Reactions) int64 {
// 	if len(*r) == 0 {
// 		return 0
// 	}
// 	// pp.Println(r)
// 	var rows int64
// 	infra.RunInTx(ctx, profile, test(r, &rows))
// 	slog.Info("Notes inserted(caller)", slog.Int64("count", rows))

// 	// return tx(ctx, infra.DB(profile), r)
// 	return rows
// }

// func tx(ctx context.Context, db bun.IDB, r *mi.Reactions) (rows int64) {
// 	// まとめて追加する(トランザクション)
// 	err := db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
// 		// JSONの中身をモデルへ移す
// 		var (
// 			users       []model.User
// 			notes       []model.Note
// 			reactions   []model.ReactionEmoji
// 			noteToTags  []model.NoteToTag
// 			files       []model.File
// 			noteToFiles []model.NoteToFile
// 		)

// 		for _, v := range *r {
// 			var dn string
// 			if v.Note.User.Name == nil {
// 				dn = v.Note.User.Username
// 			} else {
// 				dn = v.Note.User.Name.(string)
// 			}
// 			u := model.User{
// 				ID:          v.Note.User.ID,
// 				Name:        v.Note.User.Username,
// 				DisplayName: dn,
// 				AvatarURL:   v.Note.User.AvatarURL,
// 			}
// 			users = append(users, u)

// 			for _, tv := range v.Note.Tags {
// 				ht := model.HashTag{Text: tv}
// 				err := insertHashTag(ctx, tx, &ht)
// 				if err != nil {
// 					return err
// 				}
// 				slog.Info("HashTag ID", slog.Int64("ID", ht.ID))
// 				// pp.Println(ht.ID)
// 				// id is scanned automatically
// 				noteToTags = append(noteToTags, model.NoteToTag{NoteID: v.Note.ID, HashTagID: ht.ID})
// 			}

// 			for _, fv := range v.Note.Files {
// 				f := model.File{
// 					ID:           fv.ID,
// 					Name:         fv.Name,
// 					URL:          fv.URL,
// 					ThumbnailURL: fv.ThumbnailURL,
// 					Type:         fv.Type,
// 					CreatedAt:    fv.CreatedAt,
// 				}
// 				files = append(files, f)
// 				noteToFiles = append(noteToFiles, model.NoteToFile{NoteID: v.Note.ID, FileID: f.ID})
// 			}

// 			reactionName := strings.TrimSuffix(strings.TrimPrefix(v.Note.MyReaction, ":"), "@.:")
// 			n := model.Note{
// 				ID:                v.Note.ID,
// 				ReactionID:        v.ID,
// 				UserID:            v.Note.User.ID,
// 				ReactionEmojiName: reactionName,
// 				Text:              v.Note.Text,
// 				CreatedAt:         v.Note.CreatedAt, // SQLite は日時をUTCで保持する
// 			}
// 			notes = append(notes, n)

// 			r := model.ReactionEmoji{
// 				Name: reactionName,
// 			}
// 			reactions = append(reactions, r)
// 		}

// 		// 重複していたらアップデート
// 		err := insertUsers(ctx, tx, &users)
// 		if err != nil {
// 			return err
// 		}

// 		// 重複していたら登録しない(エラーにしない)
// 		rows, err = insertNotes(ctx, tx, &notes)
// 		if err != nil {
// 			return err
// 		}
// 		slog.Info("Notes inserted", slog.Int64("count", rows))

// 		err = insertReactions(ctx, tx, &reactions)
// 		if err != nil {
// 			return err
// 		}

// 		// 0件の場合がある
// 		if len(noteToTags) > 0 {
// 			err = insertNoteToTags(ctx, tx, &noteToTags)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		if len(files) > 0 {
// 			err = insertFiles(ctx, tx, &files)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		if len(noteToFiles) > 0 {
// 			err = insertNoteToFiles(ctx, tx, &noteToFiles)
// 			if err != nil {
// 				return err
// 			}
// 		}

// 		err = count(ctx, tx)
// 		// TODO: あってもなくても変わらない vs 統一感
// 		if err != nil {
// 			return err
// 		}
// 		return err
// 	})
// 	if err != nil {
// 		slog.Error(err.Error())
// 		panic(err)
// 	}
// 	return
// }

func (infra *Infra) InsertHashTag(ctx context.Context, db bun.IDB, hashtag *model.HashTag) error {
	_, err := db.NewInsert().Model(hashtag).On("CONFLICT DO UPDATE").Exec(ctx)
	return err
}

func (infra *Infra) InsertUsers(ctx context.Context, db bun.IDB, users *[]model.User) error {
	_, err := db.NewInsert().Model(users).On("CONFLICT DO UPDATE").Exec(ctx)
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

func (infra *Infra) InsertFiles(ctx context.Context, db bun.IDB, files *[]model.File) error {
	_, err := db.NewInsert().Model(files).On("CONFLICT DO UPDATE").Exec(ctx)
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

func (infra *Infra) InsertEmoji(ctx context.Context, profile string, id int64, e *mi.Emoji) {
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

func (infra *Infra) InsertColor(ctx context.Context, profile, id, c1, c2 string) {
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
