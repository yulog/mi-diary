package model

import (
	"time"
)

type Archive struct {
	YM       string
	YmCount  int64
	YMD      string
	YmdCount int64
}

type DisplayNote struct {
	ID           string `bun:",pk"`
	UserID       string `bun:",pk"`
	UserName     string
	DisplayName  string
	AvatarURL    string
	ReactionName string
	Text         string
	CreatedAt    time.Time
	User         User      `bun:"rel:belongs-to,join:user_id=id"`
	Reaction     Reaction  `bun:"rel:belongs-to,join:reaction_name=name"`
	Tags         []HashTag `bun:"m2m:note_to_tags,join:Note=HashTag"`
}
