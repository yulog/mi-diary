package model

import "github.com/uptrace/bun"

type Note struct {
	bun.BaseModel `bun:"table:notes,alias:n"`

	ID           string `bun:",pk"`
	UserID       string `bun:",pk"`
	ReactionName string
	User         User      `bun:"rel:belongs-to,join:user_id=id"`
	Reaction     Reaction  `bun:"rel:belongs-to,join:reaction_name=name"`
	Tags         []HashTag `bun:"m2m:note_to_tags,join:Note=HashTag"`
}

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID    string `bun:",pk"`
	Name  string
	Count int64
}

type Reaction struct {
	bun.BaseModel `bun:"table:reactions,alias:r"`

	Name  string `bun:",pk"`
	Image string
	Count int64
}

type HashTag struct {
	bun.BaseModel `bun:"table:hash_tags,alias:h"`

	ID    int64  `bun:",pk,autoincrement"`
	Text  string `bun:",unique"`
	Notes []Note `bun:"m2m:note_to_tags,join:HashTag=Note"`
}

type NoteToTag struct {
	NoteID    string   `bun:",pk"`
	Note      *Note    `bun:"rel:belongs-to,join:note_id=id"`
	HashTagID int64    `bun:",pk"`
	HashTag   *HashTag `bun:"rel:belongs-to,join:hash_tag_id=id"`
}
