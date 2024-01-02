package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/yulog/mi-diary/mi"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"

	"github.com/k0kubun/pp/v3"
)

// type User struct {
// 	bun.BaseModel `bun:"table:users,alias:u"`

// 	ID   int64 `bun:",pk,autoincrement"`
// 	Name string
// }

// type Book struct {
// 	bun.BaseModel `bun:"table:books,alias:b"`

// 	ID   int64 `bun:",pk"`
// 	Name string
// }

type Note struct {
	bun.BaseModel `bun:"table:notes,alias:n"`

	ID     string `bun:",pk"`
	UserID string `bun:",pk"`
	User   User   `bun:"rel:belongs-to,join:user_id=id"`
}

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID    string `bun:",pk"`
	Name  string
	Count int64
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
	_, _ = db.NewCreateTable().Model((*Note)(nil)).Exec(ctx)
	_, _ = db.NewCreateTable().Model((*User)(nil)).Exec(ctx)
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

	// Tableを作る
	// _, err = db.NewCreateTable().Model((*Book)(nil)).Exec(ctx)
	// Insert
	// 重複してたら登録しない(エラーにしない)
	// book := &Book{ID: 1, Name: "admin"}
	// _, err = db.NewInsert().Model(book).Ignore().Exec(ctx)
	// Select
	// var books []Book
	// err = db.NewSelect().Model(&books).OrderExpr("id ASC").Limit(10).Scan(ctx)
	// fmt.Println(books)

	// JSON読み込み
	f, _ := os.ReadFile("users_reactions.json")
	var r mi.Reactions
	json.Unmarshal(f, &r)
	// fmt.Printf("%+v\n", r)
	pp.Println(r)

	// JSONの中身をモデルへ移す
	var (
		users []User
		notes []Note
	)
	for _, v := range r {
		u := User{
			ID:   v.Note.User.ID,
			Name: v.Note.User.Username,
		}
		users = append(users, u)

		n := Note{
			ID:     v.Note.ID,
			UserID: v.Note.User.ID,
		}
		notes = append(notes, n)
	}

	// まとめて追加する
	_, err = db.NewInsert().Model(&users).Ignore().Exec(ctx)
	if err != nil {
		panic(err)
	}
	for _, user := range users {
		fmt.Println(user.ID) // id is scanned automatically
	}
	db.NewInsert().Model(&notes).Ignore().Exec(ctx)
	if err != nil {
		panic(err)
	}
	for _, note := range notes {
		fmt.Println(note.ID) // id is scanned automatically
	}

	// 結合
	db.NewSelect().
		Model(&notes). // &必須
		Relation("User").
		Scan(ctx)
	pp.Println(notes)

	// 特定ユーザーのノートを取得
	db.NewSelect().
		Model(&notes).
		Relation("User").
		Where("user_id = ?", "7rkrarq81i").
		Scan(ctx)
	pp.Println(notes)
}
