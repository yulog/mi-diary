package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/yulog/mi-diary/mi"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"

	"github.com/k0kubun/pp/v3"
)

type Note struct {
	bun.BaseModel `bun:"table:notes,alias:n"`

	ID           string `bun:",pk"`
	UserID       string `bun:",pk"`
	ReactionName string
	User         User      `bun:"rel:belongs-to,join:user_id=id"`
	Reaction     Reaction  `bun:"rel:belongs-to,join:reaction_name=name"`
	Tags         []HashTag `bun:"m2m:note_to_tags,join:Note=HashTag"`
}

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID    string `bun:",pk"`
	Name  string
	Count int64
}

type Reaction struct {
	bun.BaseModel `bun:"table:reactions,alias:r"`

	Name  string `bun:",pk"`
	Image string
	Count int64
}

type HashTag struct {
	bun.BaseModel `bun:"table:HashTags,alias:h"`

	ID    int64  `bun:",pk,autoincrement"`
	Text  string `bun:",unique"`
	Notes []Note `bun:"m2m:note_to_tags,join:HashTag=Note"`
}

type NoteToTag struct {
	NoteID    string   `bun:",pk"`
	Note      *Note    `bun:"rel:belongs-to,join:note_id=id"`
	HashTagID int64    `bun:",pk"`
	HashTag   *HashTag `bun:"rel:belongs-to,join:hash_tag_id=id"`
}

func main() {
	// sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	sqldb, err := sql.Open(sqliteshim.ShimName, "file:diary.db")
	if err != nil {
		panic(err)
	}
	db := bun.NewDB(sqldb, sqlitedialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))
	ctx := context.Background()
	// modelを最初に使う前にやる
	db.RegisterModel((*NoteToTag)(nil))

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
	// TODO: Migration
	_, _ = db.NewCreateTable().Model((*Note)(nil)).Exec(ctx)
	_, _ = db.NewCreateTable().Model((*User)(nil)).Exec(ctx)
	_, _ = db.NewCreateTable().Model((*Reaction)(nil)).Exec(ctx)
	_, _ = db.NewCreateTable().Model((*HashTag)(nil)).Exec(ctx)
	_, _ = db.NewCreateTable().Model((*NoteToTag)(nil)).Exec(ctx)
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
		users      []User
		notes      []Note
		reactions  []Reaction
		noteToTags []NoteToTag
	)

	// まとめて追加する(トランザクション)
	err = db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		// JSONの中身をモデルへ移す
		for _, v := range r {
			u := User{
				ID:   v.Note.User.ID,
				Name: v.Note.User.Username,
			}
			users = append(users, u)

			for _, tv := range v.Note.Tags {
				ht := HashTag{Text: tv}
				_, err = db.NewInsert().Model(&ht).On("CONFLICT DO UPDATE").Exec(ctx)
				pp.Println(ht.ID)
				noteToTags = append(noteToTags, NoteToTag{NoteID: v.Note.ID, HashTagID: ht.ID})
			}

			reactionName := strings.TrimSuffix(strings.TrimPrefix(v.Note.MyReaction, ":"), "@.:")
			n := Note{
				ID:           v.Note.ID,
				UserID:       v.Note.User.ID,
				ReactionName: reactionName,
			}
			notes = append(notes, n)

			r := Reaction{
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
	err = db.NewSelect().Model((*Note)(nil)).ColumnExpr("reaction_name as name, count(*) as count").Group("reaction_name").Scan(ctx, &reactions)
	pp.Println(reactions)

	// 既存の値も更新している？
	_, err = db.NewUpdate().
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
