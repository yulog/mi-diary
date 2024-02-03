package infra

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/model"
	"github.com/yulog/mi-diary/util/pg"
)

func (infra *Infra) Reactions(ctx context.Context, profile string) []model.Reaction {
	var reactions []model.Reaction
	infra.DB(profile).
		NewSelect().
		Model(&reactions).
		Order("count DESC").
		Scan(ctx)
	return reactions
}

func (infra *Infra) HashTags(ctx context.Context, profile string) []model.HashTag {
	var tags []model.HashTag
	infra.DB(profile).
		NewSelect().
		Model(&tags).
		Order("count DESC").
		Scan(ctx)
	return tags
}

func (infra *Infra) Users(ctx context.Context, profile string) []model.User {
	var users []model.User
	infra.DB(profile).
		NewSelect().
		Model(&users).
		Order("count DESC").
		Scan(ctx)
	return users
}

func (infra *Infra) ReactionNotes(ctx context.Context, profile, name string) []model.Note {
	var notes []model.Note
	infra.DB(profile).
		NewSelect().
		Model(&notes).
		Where("reaction_name = ?", name).
		Order("created_at DESC").
		Scan(ctx)
	return notes
}

func (infra *Infra) HashTagNotes(ctx context.Context, profile, name string) []model.Note {
	var notes []model.Note
	infra.DB(profile).
		NewSelect().
		Model((*model.NoteToTag)(nil)).
		// 必要な列だけ選択して、不要な列をなくす
		Relation("Note", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("id", "user_id", "reaction_name", "text")
		}).
		Relation("HashTag", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("")
		}).
		Column("").
		Where("hash_tag.text = ?", name).
		Order("created_at DESC").
		Scan(ctx, &notes)
	return notes
}

func (infra *Infra) UserNotes(ctx context.Context, profile, name string) []model.Note {
	var notes []model.Note
	infra.DB(profile).
		NewSelect().
		Model(&notes).
		Relation("User").
		Where("user.name = ?", name).
		Order("created_at DESC").
		Scan(ctx)
	return notes
}

func (infra *Infra) NoteCount(ctx context.Context, profile string) (int, error) {
	return infra.DB(profile).
		NewSelect().
		Model((*model.Note)(nil)).
		Count(ctx)
}

func (infra *Infra) Notes(ctx context.Context, profile string, p *pg.Pager) []model.Note {
	var notes []model.Note
	infra.DB(profile).
		NewSelect().
		Model(&notes).
		// Relation("User").
		Order("created_at DESC").
		Limit(p.Limit()).
		Offset(p.Offset()).
		Scan(ctx)
	return notes
}

func (infra *Infra) Archives(ctx context.Context, profile string) []model.Archive {
	var archives []model.Archive
	infra.DB(profile).
		NewSelect().
		Model((*model.Day)(nil)).
		Relation("Month", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Column("")
		}).
		ColumnExpr("d.ym as ym, month.count as ym_count, d.ymd as ymd, d.count as ymd_count").
		Order("ym DESC", "ymd DESC").
		Scan(ctx, &archives)
	return archives
}

func (infra *Infra) ArchiveNotes(ctx context.Context, profile, col, d string, p *pg.Pager) []model.Note {
	var notes []model.Note
	infra.DB(profile).
		NewSelect().
		Model(&notes).
		Where(col+" = ?", d). // 条件指定に関数適用した列を使う
		Order("created_at DESC").
		Limit(p.Limit()).
		Offset(p.Offset()).
		Scan(ctx)
	return notes
}
