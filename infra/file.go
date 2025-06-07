package infra

import (
	"context"
	"log/slog"

	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/domain/model"
	"github.com/yulog/mi-diary/domain/repository"
	"github.com/yulog/mi-diary/util/pagination"
)

type FileInfra struct {
	dao *DataBase
}

func (i *Infra) NewFileInfra() repository.FileRepositorier {
	return &FileInfra{dao: i.dao}
}

func (i *FileInfra) Get(ctx context.Context, profile, color string, p pagination.Paging) ([]model.File, error) {
	var files []model.File
	qb := i.dao.DB(profile).NewSelect().QueryBuilder()
	qb = addWhere(qb, "f.group_color", color)
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

func (i *FileInfra) GetByNoteID(ctx context.Context, profile, id string) ([]model.File, error) {
	// サブクエリを使う
	// note idだけ必要
	sq := i.dao.DB(profile).
		NewSelect().
		Model((*model.NoteToFile)(nil)).
		Column("file_id").
		Where("note_id = ?", id)

	var files []model.File
	err := i.dao.DB(profile).
		NewSelect().
		Model(&files).
		Where("f.id IN (?)", sq). // サブクエリを使う
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (i *FileInfra) GetByEmptyColor(ctx context.Context, profile string) ([]model.File, error) {
	var files []model.File
	err := i.dao.DB(profile).
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

func (i *FileInfra) Count(ctx context.Context, profile string) (int, error) {
	return i.dao.DB(profile).
		NewSelect().
		Model((*model.File)(nil)).
		Count(ctx)
}

func (i *FileInfra) Insert(ctx context.Context, profile string, files *[]model.File) error {
	db, ok := txFromContext(ctx)
	if !ok {
		db = i.dao.DB(profile)
	}
	_, err := db.NewInsert().
		Model(files).
		On("CONFLICT DO UPDATE").
		Exec(ctx)
	return err
}

func (i *FileInfra) UpdateByPKWithColor(ctx context.Context, profile, id, color1, color2 string) {
	r := model.File{
		ID:            id,
		DominantColor: color1,
		GroupColor:    color2,
	}
	var s []model.File
	s = append(s, r)
	_, err := i.dao.DB(profile).NewUpdate().
		Model(&s).
		OmitZero().
		Column("dominant_color").
		Column("group_color").
		Bulk().
		Exec(ctx)
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}
