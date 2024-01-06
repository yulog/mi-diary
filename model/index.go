package model

import (
	"github.com/uptrace/bun"
)

type IndexQueryCreator func(db *bun.DB) *bun.CreateIndexQuery

var IdxCreators = []IndexQueryCreator{
	// UNIQUEになっているとAtlasで自動で作られる？
	// sqlite_autoindexというのも作られるっぽい？
	// func(db *bun.DB) *bun.CreateIndexQuery {
	// 	return db.NewCreateIndex().
	// 		Model((*HashTag)(nil)).
	// 		Index("hash_tags_text").
	// 		Column("text")
	// },
}
