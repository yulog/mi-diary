package app

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-json"
	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/mi"
	"github.com/yulog/mi-diary/model"
)

func (app *App) InsertFromFile(ctx context.Context, profile string) {
	// JSON読み込み
	f, _ := os.ReadFile("users_reactions.json")

	app.Insert(ctx, profile, f)
}

func (app *App) Insert(ctx context.Context, profile string, b []byte) {
	// app := New()
	db := app.DB(profile)
	// JSON Unmarshal
	var r mi.Reactions
	json.Unmarshal(b, &r)
	// pp.Println(r)

	tx(ctx, db, r)
}

func tx(ctx context.Context, db *bun.DB, r mi.Reactions) {
	// JSONの中身をモデルへ移す
	var (
		users      []model.User
		notes      []model.Note
		reactions  []model.Reaction
		noteToTags []model.NoteToTag
	)

	// まとめて追加する(トランザクション)
	err := db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		// JSONの中身をモデルへ移す
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
			}
			users = append(users, u)

			for _, tv := range v.Note.Tags {
				ht := model.HashTag{Text: tv}
				_, _ = db.NewInsert().Model(&ht).On("CONFLICT DO UPDATE").Exec(ctx)
				// pp.Println(ht.ID)
				noteToTags = append(noteToTags, model.NoteToTag{NoteID: v.Note.ID, HashTagID: ht.ID})
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

		// 重複してたら登録しない(エラーにしない)
		_, err := db.NewInsert().Model(&users).Ignore().Exec(ctx)
		if err != nil {
			return err
		}
		// for _, user := range users {
		// 	fmt.Println(user.ID) // id is scanned automatically
		// }

		_, err = db.NewInsert().Model(&notes).Ignore().Exec(ctx)
		if err != nil {
			return err
		}
		// for _, note := range notes {
		// 	fmt.Println(note.ID) // id is scanned automatically
		// }

		_, err = db.NewInsert().Model(&reactions).Ignore().Exec(ctx)
		if err != nil {
			return err
		}
		// for _, reaction := range reactions {
		// 	fmt.Println(reaction.ID) // id is scanned automatically
		// }

		// 0件の場合がある
		if len(noteToTags) > 0 {
			_, err = db.NewInsert().Model(&noteToTags).Ignore().Exec(ctx)
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
		ColumnExpr("reaction_name as name, count(*) as count").
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
		ColumnExpr("hash_tag_id as id, count(*) as count").
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
		ColumnExpr("strftime('%Y-%m', created_at, 'localtime') as ym, count(*) as count").
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
		ColumnExpr("strftime('%Y-%m-%d', created_at, 'localtime') as ymd, strftime('%Y-%m', created_at, 'localtime') as ym, count(*) as count").
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

func (app *App) InsertEmoji(ctx context.Context, profile string, b []byte) {
	// app := New()
	db := app.DB(profile)
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
