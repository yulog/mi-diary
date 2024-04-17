package migrate

import (
	"os"

	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/app"
	"github.com/yulog/mi-diary/infra"
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
	infra := infra.New(app)
	models := []interface{}{
		(*model.Note)(nil),
		(*model.User)(nil),
		(*model.Reaction)(nil),
		(*model.HashTag)(nil),
		(*model.NoteToTag)(nil),
		(*model.File)(nil),
		(*model.NoteToFile)(nil),
		(*model.Month)(nil),
		(*model.Day)(nil),
	}
	var data []byte
	for k := range app.Config.Profiles {
		data = append(data, modelsToByte(infra.DB(k), models)...)
		data = append(data, indexesToByte(infra.DB(k), model.IdxCreators)...)
		break // schemaの生成は1つだけやれば良さそう
	}
	// TODO: 権限これで良いの？
	os.WriteFile("migrate/schema.sql", data, 0666)
}
