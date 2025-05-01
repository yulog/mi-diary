package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"net/http"
	"net/url"
	"os"
	"strings"

	_ "golang.org/x/image/webp"

	icolor "github.com/yulog/mi-diary/color"
	"github.com/yulog/mi-diary/domain/model"
	mi "github.com/yulog/miutil"
	"github.com/yulog/miutil/miauth"

	"github.com/uptrace/bun"

	"github.com/k0kubun/pp/v3"

	"github.com/cenkalti/dominantcolor"
	"github.com/mattn/go-ciede2000"
)

func main() {
	resp, _ := http.Get("")
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	red := &color.RGBA{204, 0, 0, 255}
	green := &color.RGBA{0, 204, 0, 255}
	bluegreen := &color.RGBA{3, 192, 198, 255}
	blue := &color.RGBA{0, 0, 255, 255}
	diff := ciede2000.Diff(dominantcolor.Find(img), green)
	fmt.Println(diff)
	fmt.Println(ciede2000.Diff(dominantcolor.Find(img), red))
	fmt.Println(ciede2000.Diff(dominantcolor.Find(img), bluegreen))
	fmt.Println(ciede2000.Diff(dominantcolor.Find(img), blue))
	fmt.Println(dominantcolor.Hex(dominantcolor.Find(img)))
	hex, hex2, err := icolor.Color("")
	fmt.Println(hex, hex2)
}

func miauthexp() {
	u, _ := url.Parse("https://misskey.io")
	conf := &miauth.AuthConfig{
		Name:       "mi-diary-test",
		Callback:   fmt.Sprintf("http://localhost:1323/callback/%s", u.Host),
		Permission: []string{"read:reactions"},
		Host:       "https://misskey.io",
	}
	fmt.Println(conf.AuthCodeURL())
}

func createTable() {
	// app := app.New()
	// infra := infra.New(app)
	// db := infra.DB("")
	db := bun.DB{}

	ctx := context.Background()

	// 生のSQL
	// res, err := db.ExecContext(ctx, "SELECT 1")
	// fmt.Println(res.RowsAffected())
	// var num int
	// err = db.QueryRowContext(ctx, "SELECT 1").Scan(&num)
	// fmt.Println(num)

	// Bunのクエリビルダー
	// res, err := db.NewSelect().ColumnExpr("1").Exec(ctx)
	// fmt.Println(res.LastInsertId())
	// var num int
	// err = db.NewSelect().ColumnExpr("1").Scan(ctx, &num)
	// fmt.Println(num)

	// Tableを作る
	_, _ = db.NewCreateTable().Model((*model.Note)(nil)).Exec(ctx)
	_, _ = db.NewCreateTable().Model((*model.User)(nil)).Exec(ctx)
	_, _ = db.NewCreateTable().Model((*model.ReactionEmoji)(nil)).Exec(ctx)
	_, _ = db.NewCreateTable().Model((*model.HashTag)(nil)).Exec(ctx)
	_, _ = db.NewCreateTable().Model((*model.NoteToTag)(nil)).Exec(ctx)
	// Insert
	// user := &User{Name: "admin"}
	// _, err = db.NewInsert().Model(user).Exec(ctx)

	// Update
	// userupd := &User{ID: 2, Name: "user2"}
	// _, err = db.NewUpdate().Model(userupd).Column("name").WherePK().Exec(ctx)

	// Select
	// var users []User
	// err = db.NewSelect().Model(&users).OrderExpr("id ASC").Limit(10).Scan(ctx)
	// fmt.Println(users)

	// JSON読み込み
	f, _ := os.ReadFile("users_reactions.json")
	var r mi.Reactions
	json.Unmarshal(f, &r)
	// pp.Println(r)

	// JSONの中身をモデルへ移す
	var (
		users      []model.User
		notes      []model.Note
		reactions  []model.ReactionEmoji
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
				pp.Println(ht.ID)
				noteToTags = append(noteToTags, model.NoteToTag{NoteID: v.Note.ID, HashTagID: ht.ID})
			}

			reactionName := strings.TrimSuffix(strings.TrimPrefix(v.Note.MyReaction, ":"), "@.:")
			n := model.Note{
				ID:                v.Note.ID,
				UserID:            v.Note.User.ID,
				ReactionEmojiName: reactionName,
			}
			notes = append(notes, n)

			r := model.ReactionEmoji{
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

	// 結合
	db.NewSelect().
		Model(&notes). // &必須
		Relation("User").
		Relation("Reaction").
		Scan(ctx)
	pp.Println(notes)

	// 特定ユーザーのノートを取得
	db.NewSelect().
		Model(&notes).
		Relation("User").
		Relation("Reaction").
		Where("user_id = ?", "7rkrarq81i").
		Scan(ctx)
	pp.Println(notes)

	// リアクション
	db.NewSelect().
		Model(&reactions).
		Scan(ctx)
	pp.Println(reactions)

	// var m []map[string]interface{}
	// err = db.NewSelect().Model((*Note)(nil)).ColumnExpr("reaction_name, count(*)").Group("reaction_name").Scan(ctx, &m)
	// pp.Println(m)

	// 先に取得してあったreactionsを更新している？
	_ = db.NewSelect().Model((*model.Note)(nil)).ColumnExpr("reaction_name as name, count(*) as count").Group("reaction_name").Scan(ctx, &reactions)
	pp.Println(reactions)

	// 既存の値も更新している？
	_, _ = db.NewUpdate().
		Model(&reactions).
		OmitZero().
		Column("count").
		Bulk().
		Exec(ctx)
	// リアクション
	db.NewSelect().
		Model(&reactions).
		Scan(ctx)
	pp.Println(reactions)

	// 特定Tagのノートを取得
	db.NewSelect().
		Model(&noteToTags).
		Relation("HashTag").
		Relation("Note").
		Where("hash_tag_id = ?", 1).
		Scan(ctx)
	pp.Println(noteToTags)
}
