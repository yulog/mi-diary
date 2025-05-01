package repository

import mi "github.com/yulog/miutil"

type MisskeyAPIRepositorier interface {
	GetUserReactions(profile, id string, limit int) (int, *mi.Reactions, error)
	GetEmoji(profile, name string) (*mi.Emoji, error)
}
