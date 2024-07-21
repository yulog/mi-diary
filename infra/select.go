package infra

import (
	"context"

	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/model"
	"github.com/yulog/mi-diary/util/pg"
)

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

func (infra *Infra) Files(ctx context.Context, profile string, p *pg.Pager) ([]model.File, error) {
	var files []model.File
	err := infra.DB(profile).
		NewSelect().
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

func (infra *Infra) NoteCount(ctx context.Context, profile string) (int, error) {
	return infra.DB(profile).
		NewSelect().
		Model((*model.Note)(nil)).
		Count(ctx)
}

func (infra *Infra) Notes(ctx context.Context, profile, s string, p *pg.Pager) ([]model.Note, error) {
	var notes []model.Note
	q := infra.DB(profile).
		NewSelect().
		Model(&notes).
		Relation("User").
		Relation("Files").
		Order("created_at DESC").
		Limit(p.Limit()).
		Offset(p.Offset())

	if s != "" {
		q.Where("text like ?", "%"+s+"%")
	}
	err := q.Scan(ctx)
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

func (infra *Infra) ArchiveNotes(ctx context.Context, profile, col, d string, p *pg.Pager) ([]model.Note, error) {
	var notes []model.Note
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
