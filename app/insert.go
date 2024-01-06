package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/mi"
	"github.com/yulog/mi-diary/model"
)

func Insert(ctx context.Context) {
	app := New()
	db := app.DB()
	// JSON読み込み
	f, _ := os.ReadFile("users_reactions.json")
	var r mi.Reactions
	json.Unmarshal(f, &r)
	// pp.Println(r)

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
			u := model.User{
				ID:   v.Note.User.ID,
				Name: v.Note.User.Username,
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
			}
			notes = append(notes, n)

			r := model.Reaction{
				Name:  reactionName,
				Image: "xxx",
			}
			reactions = append(reactions, r)
		}

		// 重複してたら登録しない(エラーにしない)
		_, err := db.NewInsert().Model(&users).Ignore().Exec(ctx)
		if err != nil {
			return err
		}
		for _, user := range users {
			fmt.Println(user.ID) // id is scanned automatically
		}

		_, err = db.NewInsert().Model(&notes).Ignore().Exec(ctx)
		if err != nil {
			return err
		}
		for _, note := range notes {
			fmt.Println(note.ID) // id is scanned automatically
		}

		_, err = db.NewInsert().Model(&reactions).Ignore().Exec(ctx)
		if err != nil {
			return err
		}
		// for _, reaction := range reactions {
		// 	fmt.Println(reaction.ID) // id is scanned automatically
		// }

		_, err = db.NewInsert().Model(&noteToTags).Ignore().Exec(ctx)
		if err != nil {
			return err
		}
		return err
	})
	if err != nil {
		panic(err)
	}
}
