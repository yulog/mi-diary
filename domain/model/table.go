package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Note struct {
	bun.BaseModel `bun:"table:notes,alias:n"`

	ID                string `bun:",pk"`
	ReactionID        string `bun:",notnull"`
	UserID            string `bun:",notnull"` // `bun:",pk"` ここをprimary keyにするとm2mのリレーション結合が壊れる
	ReactionEmojiName string `bun:",notnull"`
	Text              string
	CreatedAt         time.Time
	User              User          `bun:"rel:belongs-to,join:user_id=id"`
	Reaction          ReactionEmoji `bun:"rel:belongs-to,join:reaction_emoji_name=name"`
	Tags              []HashTag     `bun:"m2m:note_to_tags,join:Note=HashTag"`
	Files             []File        `bun:"m2m:note_to_files,join:Note=File"`
}

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID          string `bun:",pk"`
	Name        string
	DisplayName string
	AvatarURL   string
	Count       int64
}

type ReactionEmoji struct {
	bun.BaseModel `bun:"table:reactions,alias:r"`

	ID       int64  `bun:",pk,autoincrement"`
	Name     string `bun:",unique,notnull"`
	Image    string
	IsSymbol bool
	Count    int64
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

type File struct {
	bun.BaseModel `bun:"table:files,alias:f"`

	ID            string `bun:",pk"`
	Name          string
	URL           string
	ThumbnailURL  string
	Type          string
	DominantColor string
	GroupColor    string
	CreatedAt     time.Time
	Notes         []Note `bun:"m2m:note_to_files,join:File=Note"`
}

type NoteToFile struct {
	NoteID string `bun:",pk"`
	Note   *Note  `bun:"rel:belongs-to,join:note_id=id"`
	FileID string `bun:",pk"`
	File   *File  `bun:"rel:belongs-to,join:file_id=id"`
}

type Month struct {
	bun.BaseModel `bun:"table:months,alias:m"`

	YM    string `bun:",pk"`
	Count int64
	Days  []*Day `bun:"rel:has-many,join:ym=ym"`
}

type Day struct {
	bun.BaseModel `bun:"table:days,alias:d"`

	YMD   string `bun:",pk"`
	YM    string
	Count int64
	Month Month `bun:"rel:belongs-to,join:ym=ym"`
}
