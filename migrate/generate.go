package migrate

import (
	"os"

	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/domain/model"
)

// https://techblog.enechain.com/entry/bun-atlas-migration-setup-guide
func modelsToByte(db *bun.DB, models []any) []byte {
	var data []byte
	for _, model := range models {
		query := db.NewCreateTable().Model(model).WithForeignKeys()
		rawQuery, err := query.AppendQuery(db.Formatter(), nil)
		if err != nil {
			panic(err)
		}
		data = append(data, rawQuery...)
		data = append(data, ";\n"...)
	}
	return data
}

// https://techblog.enechain.com/entry/bun-atlas-migration-setup-guide
func indexesToByte(db *bun.DB, idxCreators []model.IndexQueryCreator) []byte {
	var data []byte
	for _, idxCreator := range idxCreators {
		idx := idxCreator(db)
		rawQuery, err := idx.AppendQuery(db.Formatter(), nil)
		if err != nil {
			panic(err)
		}
		data = append(data, rawQuery...)
		data = append(data, ";\n"...)
	}
	return data
}

func GenerateSchema(db *bun.DB) {
	models := []any{
		(*model.Note)(nil),
		(*model.User)(nil),
		(*model.ReactionEmoji)(nil),
		(*model.HashTag)(nil),
		(*model.NoteToTag)(nil),
		(*model.File)(nil),
		(*model.NoteToFile)(nil),
		(*model.Month)(nil),
		(*model.Day)(nil),
	}
	var data []byte
	data = append(data, modelsToByte(db, models)...)
	data = append(data, indexesToByte(db, model.IdxCreators)...)
	// TODO: 権限これで良いの？
	os.WriteFile("migrate/schema.sql", data, 0666)
}
