package infra

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/logic"
	"github.com/yulog/mi-diary/model"
	"github.com/yulog/mi-diary/util/pg"
)

type FileInfra struct {
	infra *Infra
}

func (i *Infra) NewFileInfra() logic.FileRepositorier {
	return &FileInfra{infra: i}
}

func (fi *FileInfra) Get(ctx context.Context, profile, c string, p *pg.Pager) ([]model.File, error) {
	var files []model.File
	qb := fi.infra.DB(profile).NewSelect().QueryBuilder()
	qb = addWhere(qb, "f.group_color", c)
	err := qb.
		Unwrap().(*bun.SelectQuery).
		Model(&files).
		Relation("Notes", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Relation("User")
		}).
		Order("created_at DESC").
		Limit(p.Limit()).
		Offset(p.Offset()).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (fi *FileInfra) GetByNoteID(ctx context.Context, profile, id string) ([]model.File, error) {
	// サブクエリを使う
	// note idだけ必要
	sq := fi.infra.DB(profile).
		NewSelect().
		Model((*model.NoteToFile)(nil)).
		Column("file_id").
		Where("note_id = ?", id)

	var files []model.File
	err := fi.infra.DB(profile).
		NewSelect().
		Model(&files).
		Where("f.id IN (?)", sq). // サブクエリを使う
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (fi *FileInfra) GetByEmptyColor(ctx context.Context, profile string) ([]model.File, error) {
	var files []model.File
	err := fi.infra.DB(profile).
		NewSelect().
		Model(&files).
		Where("group_color = ?", "").
		WhereOr("group_color is null").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (fi *FileInfra) Count(ctx context.Context, profile string) (int, error) {
	return fi.infra.DB(profile).
		NewSelect().
		Model((*model.File)(nil)).
		Count(ctx)
}

func (fi *FileInfra) Insert(ctx context.Context, db bun.IDB, files *[]model.File) error {
	_, err := db.NewInsert().
		Model(files).
		On("CONFLICT DO UPDATE").
		Exec(ctx)
	return err
}
