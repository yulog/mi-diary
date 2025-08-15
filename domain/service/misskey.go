package service

import (
	mi "github.com/yulog/miutil"
)

type MisskeyAPIServicer interface {
	Client(host, credential string) MisskeyAPIServicer

	GetUserReactions(userID, untilID string, limit int) (int, *mi.Reactions, error)
	GetEmoji(name string) (*mi.Emoji, error)
	GetEmojis() (*mi.Emojis, error)
}
