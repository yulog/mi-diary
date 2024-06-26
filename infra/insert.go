package infra

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-json"
	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/model"
	mi "github.com/yulog/miutil"
)

func (infra *Infra) InsertFromFile(ctx context.Context, profile string) {
	// JSON読み込み
	f, _ := os.ReadFile("users_reactions.json")
	var r mi.Reactions
	json.Unmarshal(f, &r)

	infra.Insert(ctx, profile, &r)
}

func (infra *Infra) Insert(ctx context.Context, profile string, r *mi.Reactions) int64 {
	if len(*r) == 0 {
		return 0
	}
	// pp.Println(r)

	return tx(ctx, infra.DB(profile), r)
}

func tx(ctx context.Context, db *bun.DB, r *mi.Reactions) (rows int64) {
	// まとめて追加する(トランザクション)
	err := db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		// JSONの中身をモデルへ移す
		var (
			users       []model.User
			notes       []model.Note
			reactions   []model.ReactionEmoji
			noteToTags  []model.NoteToTag
			files       []model.File
			noteToFiles []model.NoteToFile
		)

		for _, v := range *r {
			var dn string
			if v.Note.User.Name == nil {
				dn = v.Note.User.Username
			} else {
				dn = v.Note.User.Name.(string)
			}
			u := model.User{
				ID:          v.Note.User.ID,
				Name:        v.Note.User.Username,
				DisplayName: dn,
				AvatarURL:   v.Note.User.AvatarURL,
			}
			users = append(users, u)

			for _, tv := range v.Note.Tags {
				ht := model.HashTag{Text: tv}
				_, _ = tx.NewInsert().Model(&ht).On("CONFLICT DO UPDATE").Exec(ctx)
				// pp.Println(ht.ID)
				// id is scanned automatically
				noteToTags = append(noteToTags, model.NoteToTag{NoteID: v.Note.ID, HashTagID: ht.ID})
			}

			for _, fv := range v.Note.Files {
				f := model.File{
					ID:           fv.ID,
					Name:         fv.Name,
					URL:          fv.URL,
					ThumbnailURL: fv.ThumbnailURL,
					Type:         fv.Type,
					CreatedAt:    fv.CreatedAt,
				}
				files = append(files, f)
				noteToFiles = append(noteToFiles, model.NoteToFile{NoteID: v.Note.ID, FileID: f.ID})
			}

			reactionName := strings.TrimSuffix(strings.TrimPrefix(v.Note.MyReaction, ":"), "@.:")
			n := model.Note{
				ID:                v.Note.ID,
				ReactionID:        v.ID,
				UserID:            v.Note.User.ID,
				ReactionEmojiName: reactionName,
				Text:              v.Note.Text,
				CreatedAt:         v.Note.CreatedAt, // SQLite は日時をUTCで保持する
			}
			notes = append(notes, n)

			r := model.ReactionEmoji{
				Name: reactionName,
				// Image: "",
			}
			reactions = append(reactions, r)
		}

		// 重複していたらアップデート
		_, err := tx.NewInsert().Model(&users).On("CONFLICT DO UPDATE").Exec(ctx)
		if err != nil {
			return err
		}

		// 重複していたら登録しない(エラーにしない)
		result, err := tx.NewInsert().Model(&notes).Ignore().Exec(ctx)
		if err != nil {
			return err
		}
		rows, _ = result.RowsAffected()
		fmt.Println("insert:", rows)

		_, err = tx.NewInsert().Model(&reactions).Ignore().Exec(ctx)
		if err != nil {
			return err
		}

		// 0件の場合がある
		if len(noteToTags) > 0 {
			_, err = tx.NewInsert().Model(&noteToTags).Ignore().Exec(ctx)
			if err != nil {
				return err
			}
		}

		if len(files) > 0 {
			_, err = tx.NewInsert().Model(&files).On("CONFLICT DO UPDATE").Exec(ctx)
			if err != nil {
				return err
			}
		}

		if len(noteToFiles) > 0 {
			_, err = tx.NewInsert().Model(&noteToFiles).Ignore().Exec(ctx)
			if err != nil {
				return err
			}
		}

		err = count(ctx, tx)
		// TODO: あってもなくても変わらない vs 統一感
		if err != nil {
			return err
		}
		return err
	})
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return
}

func count(ctx context.Context, db bun.IDB) error {
	// リアクションのカウント
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

	// タグのカウント
	var hashtags []model.HashTag
	err = db.NewSelect().
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

	// ユーザーのカウント
	var users []model.User
	err = db.NewSelect().
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

	// 月別のカウント
	var months []model.Month
	err = db.NewSelect().
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

	// 日別のカウント
	var days []model.Day
	err = db.NewSelect().
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
	return err
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
		fmt.Println(err)
		panic(err)
	}
}
