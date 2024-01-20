package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Note struct {
	bun.BaseModel `bun:"table:notes,alias:n"`

	ID           string `bun:",pk"`
	UserID       string `bun:",pk"`
	ReactionName string
	Text         string
	CreatedAt    time.Time
	User         User      `bun:"rel:belongs-to,join:user_id=id"`
	Reaction     Reaction  `bun:"rel:belongs-to,join:reaction_name=name"`
	Tags         []HashTag `bun:"m2m:note_to_tags,join:Note=HashTag"`
}

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID          string `bun:",pk"`
	Name        string
	DisplayName string
	Count       int64
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
	Count int64
	Notes []Note `bun:"m2m:note_to_tags,join:HashTag=Note"`
}

type NoteToTag struct {
	NoteID    string   `bun:",pk"`
	Note      *Note    `bun:"rel:belongs-to,join:note_id=id"`
	HashTagID int64    `bun:",pk"`
	HashTag   *HashTag `bun:"rel:belongs-to,join:hash_tag_id=id"`
}

type Month struct {
	bun.BaseModel `bun:"table:months,alias:m"`

	YM    string `bun:",pk"`
	Count int64
}

type Day struct {
	bun.BaseModel `bun:"table:days,alias:d"`

	YMD   string `bun:",pk"`
	YM    string
	Count int64
	Month Month `bun:"rel:belongs-to,join:ym=ym"`
}

type Archive struct {
	bun.BaseModel `bun:"table:archives,alias:a"`

	YM       string
	YmCount  int64
	YMD      string
	YmdCount int64
}
