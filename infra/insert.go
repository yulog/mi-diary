package infra

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/model"
	mi "github.com/yulog/miutil"
)

func (infra *Infra) InsertFromFile(ctx context.Context, profile string) {
	// JSON読み込み
	f, _ := os.ReadFile("users_reactions.json")

	infra.Insert(ctx, profile, f)
}

func (infra *Infra) Insert(ctx context.Context, profile string, b []byte) {
	db := infra.DB(profile)
	// JSON Unmarshal
	var r mi.Reactions
	json.Unmarshal(b, &r)
	// pp.Println(r)

	tx(ctx, db, r)
}

func tx(ctx context.Context, db *bun.DB, r mi.Reactions) {
	// まとめて追加する(トランザクション)
	err := db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		// JSONの中身をモデルへ移す
		var (
			users       []model.User
			notes       []model.Note
			reactions   []model.Reaction
			noteToTags  []model.NoteToTag
			noteToFiles []model.NoteToFile
		)

		for _, v := range r {
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
				_, _ = db.NewInsert().Model(&ht).On("CONFLICT DO UPDATE").Exec(ctx)
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
				}
				_, _ = db.NewInsert().Model(&f).On("CONFLICT DO UPDATE").Exec(ctx)
				noteToFiles = append(noteToFiles, model.NoteToFile{NoteID: v.Note.ID, FileID: f.ID})
			}

			reactionName := strings.TrimSuffix(strings.TrimPrefix(v.Note.MyReaction, ":"), "@.:")
			n := model.Note{
				ID:           v.Note.ID,
				UserID:       v.Note.User.ID,
				ReactionName: reactionName,
				Text:         v.Note.Text,
				CreatedAt:    v.Note.CreatedAt, // SQLite は日時をUTCで保持する
			}
			notes = append(notes, n)

			r := model.Reaction{
				Name:  reactionName,
				Image: "",
			}
			reactions = append(reactions, r)
		}

		// 重複していたらアップデート
		_, err := db.NewInsert().Model(&users).On("CONFLICT DO UPDATE").Exec(ctx)
		if err != nil {
			return err
		}

		// 重複していたら登録しない(エラーにしない)
		result, err := db.NewInsert().Model(&notes).Ignore().Exec(ctx)
		if err != nil {
			return err
		}
		rows, _ := result.RowsAffected()
		fmt.Println("insert:", rows)
		// TODO: すべて取得するようにする際はinsert件数が0まで
		// until(?)とかを付けて繰り返す？

		_, err = db.NewInsert().Model(&reactions).Ignore().Exec(ctx)
		if err != nil {
			return err
		}

		// 0件の場合がある
		if len(noteToTags) > 0 {
			_, err = db.NewInsert().Model(&noteToTags).Ignore().Exec(ctx)
			if err != nil {
				return err
			}
		}

		if len(noteToFiles) > 0 {
			_, err = db.NewInsert().Model(&noteToFiles).Ignore().Exec(ctx)
			if err != nil {
				return err
			}
		}

		err = count(ctx, db)
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
}

func count(ctx context.Context, db *bun.DB) error {
	// リアクションのカウント
	var reactions []model.Reaction
	err := db.NewSelect().
		Model((*model.Note)(nil)).
		ColumnExpr("reaction_name as name").
		ColumnExpr("count(*) as count").
		Group("reaction_name").
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

	_, err = db.NewUpdate().
		Model(&hashtags).
		OmitZero().
		Column("count").
		Bulk().
		Exec(ctx)
	if err != nil {
		return err
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

func (infra *Infra) InsertEmoji(ctx context.Context, profile string, b []byte) {
	db := infra.DB(profile)
	// JSON Unmarshal
	var e mi.Emoji
	json.Unmarshal(b, &e)
	// pp.Println(r)

	// TODO: emoji画像をローカルに保存する

	r := model.Reaction{
		Name:  e.Name,
		Image: e.URL,
	}
	var s []model.Reaction
	s = append(s, r)
	_, err := db.NewUpdate().
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
