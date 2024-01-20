package migrate

import (
	"os"

	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/app"
	"github.com/yulog/mi-diary/model"
)

// https://techblog.enechain.com/entry/bun-atlas-migration-setup-guide
func modelsToByte(db *bun.DB, models []interface{}) []byte {
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

func GenerateSchema() {
	app := app.New()
	models := []interface{}{
		(*model.Note)(nil),
		(*model.User)(nil),
		(*model.Reaction)(nil),
		(*model.HashTag)(nil),
		(*model.NoteToTag)(nil),
		(*model.Month)(nil),
		(*model.Day)(nil),
	}
	var data []byte
	data = append(data, modelsToByte(app.DB(), models)...)
	data = append(data, indexesToByte(app.DB(), model.IdxCreators)...)
	// TODO: 権限これで良いの？
	os.WriteFile("migrate/schema.sql", data, 0666)
}
