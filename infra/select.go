package infra

import (
	"context"
	"regexp"

	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/model"
	"github.com/yulog/mi-diary/util/pg"
)

// addWhereLike
//
// https://bun.uptrace.dev/guide/query-where.html#querybuilder
func addWhereLike(q bun.QueryBuilder, col, s string) bun.QueryBuilder {
	if s == "" {
		return q
	}
	return q.Where("? like ?", bun.Ident(col), "%"+s+"%")
}

func addWhere(q bun.QueryBuilder, col, s string) bun.QueryBuilder {
	if s == "" {
		return q
	}
	return q.Where("? = ?", bun.Ident(col), s)
}

func (infra *Infra) Reactions(ctx context.Context, profile string) ([]model.ReactionEmoji, error) {
	var reactions []model.ReactionEmoji
	err := infra.DB(profile).
		NewSelect().
		Model(&reactions).
		Order("count DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return reactions, nil
}

func (infra *Infra) ReactionOne(ctx context.Context, profile, name string) (model.ReactionEmoji, error) {
	var reaction model.ReactionEmoji
	err := infra.DB(profile).
		NewSelect().
		Model(&reaction).
		Where("name = ?", name).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return model.ReactionEmoji{}, err
	}
	return reaction, nil
}

func (infra *Infra) ReactionImageEmpty(ctx context.Context, profile string) ([]model.ReactionEmoji, error) {
	var reactions []model.ReactionEmoji
	err := infra.DB(profile).
		NewSelect().
		Model(&reactions).
		Where("image = ?", "").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return reactions, nil
}

func (infra *Infra) HashTags(ctx context.Context, profile string) ([]model.HashTag, error) {
	var tags []model.HashTag
	err := infra.DB(profile).
		NewSelect().
		Model(&tags).
		Order("count DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func (infra *Infra) Users(ctx context.Context, profile string) ([]model.User, error) {
	var users []model.User
	err := infra.DB(profile).
		NewSelect().
		Model(&users).
		Order("count DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (infra *Infra) ReactionNotes(ctx context.Context, profile, name string, p *pg.Pager) ([]model.Note, error) {
	var notes []model.Note
	err := infra.DB(profile).
		NewSelect().
		Model(&notes).
		Relation("User").
		Relation("Files").
		Where("reaction_emoji_name = ?", name).
		Order("created_at DESC").
		Limit(p.Limit()).
		Offset(p.Offset()).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (infra *Infra) HashTagNotes(ctx context.Context, profile, name string, p *pg.Pager) ([]model.Note, error) {
	// サブクエリを使う
	// note idだけ必要
	sq := infra.DB(profile).
		NewSelect().
		Model((*model.NoteToTag)(nil)).
		// 必要な列だけ選択して、不要な列をなくす
		Relation("HashTag", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.ExcludeColumn("*")
		}).
		Column("note_id").
		Where("hash_tag.text = ?", name).
		Limit(p.Limit()).
		Offset(p.Offset())

	var notes []model.Note
	err := infra.DB(profile).
		NewSelect().
		Model(&notes).
		Relation("User").
		Relation("Files").
		Where("n.id IN (?)", sq). // サブクエリを使う
		Order("created_at DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (infra *Infra) UserNotes(ctx context.Context, profile, name string, p *pg.Pager) ([]model.Note, error) {
	var notes []model.Note
	err := infra.DB(profile).
		NewSelect().
		Model(&notes).
		Relation("User").
		Relation("Files").
		Where("user.name = ?", name).
		Order("created_at DESC").
		Limit(p.Limit()).
		Offset(p.Offset()).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (infra *Infra) FileCount(ctx context.Context, profile string) (int, error) {
	return infra.DB(profile).
		NewSelect().
		Model((*model.File)(nil)).
		Count(ctx)
}

func (infra *Infra) Files(ctx context.Context, profile, c string, p *pg.Pager) ([]model.File, error) {
	var files []model.File
	qb := infra.DB(profile).NewSelect().QueryBuilder()
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

func (infra *Infra) FilesByNoteID(ctx context.Context, profile, id string) ([]model.File, error) {
	// サブクエリを使う
	// note idだけ必要
	sq := infra.DB(profile).
		NewSelect().
		Model((*model.NoteToFile)(nil)).
		Column("file_id").
		Where("note_id = ?", id)

	var files []model.File
	err := infra.DB(profile).
		NewSelect().
		Model(&files).
		Where("f.id IN (?)", sq). // サブクエリを使う
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (infra *Infra) NoteCount(ctx context.Context, profile string) (int, error) {
	return infra.DB(profile).
		NewSelect().
		Model((*model.Note)(nil)).
		Count(ctx)
}

func (infra *Infra) Notes(ctx context.Context, profile, s string, p *pg.Pager) ([]model.Note, error) {
	var notes []model.Note
	// https://bun.uptrace.dev/guide/query-where.html#querybuilder
	qb := infra.DB(profile).NewSelect().QueryBuilder()
	qb = addWhereLike(qb, "text", s)
	err := qb.
		Unwrap().(*bun.SelectQuery).
		Model(&notes).
		Relation("User").
		Relation("Files").
		Order("created_at DESC").
		Limit(p.Limit()).
		Offset(p.Offset()).
		Scan(ctx)

	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (infra *Infra) Archives(ctx context.Context, profile string) ([]model.Month, error) {
	var archives []model.Month
	err := infra.DB(profile).
		NewSelect().
		Model(&archives).
		Relation("Days", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Order("ymd DESC")
		}).
		Order("ym DESC").
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return archives, nil
}

var reym = regexp.MustCompile(`^\d{4}-\d{2}$`)
var reymd = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func (infra *Infra) ArchiveNotes(ctx context.Context, profile, d string, p *pg.Pager) ([]model.Note, error) {
	var notes []model.Note
	col := ""
	if reym.MatchString(d) {
		col = "strftime('%Y-%m', created_at, 'localtime')"
	} else if reymd.MatchString(d) {
		col = "strftime('%Y-%m-%d', created_at, 'localtime')"
	}
	err := infra.DB(profile).
		NewSelect().
		Model(&notes).
		Relation("User").
		Relation("Files").
		Where("? = ?", bun.Safe(col), d). // 条件指定に関数適用した列を使う
		Order("created_at DESC").
		Limit(p.Limit()).
		Offset(p.Offset()).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return notes, nil
}
