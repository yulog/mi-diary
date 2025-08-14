package model

import (
	"github.com/uptrace/bun"
)

type IndexQueryCreator func(db *bun.DB) *bun.CreateIndexQuery

// https://techblog.enechain.com/entry/bun-atlas-migration-setup-guide
var IdxCreators = []IndexQueryCreator{
	// UNIQUEになっているとAtlasで自動で作られる？
	// sqlite_autoindexというのも作られるっぽい？
	// func(db *bun.DB) *bun.CreateIndexQuery {
	// 	return db.NewCreateIndex().
	// 		Model((*Note)(nil)).
	// 		Index("note_id").
	// 		Column("note_id")
	// },
	// func(db *bun.DB) *bun.CreateIndexQuery {
	// 	return db.NewCreateIndex().
	// 		Model((*User)(nil)).
	// 		Index("user_id").
	// 		Column("user_id")
	// },
	// func(db *bun.DB) *bun.CreateIndexQuery {
	// 	return db.NewCreateIndex().
	// 		Model((*File)(nil)).
	// 		Index("file_id").
	// 		Column("file_id")
	// },
}
