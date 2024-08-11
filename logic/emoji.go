package logic

import (
	"bytes"
	"context"
	"log/slog"
	"net/url"

	"github.com/goccy/go-json"
	"github.com/yulog/mi-diary/infra"
	mi "github.com/yulog/miutil"
)

type EmojiLogic interface {
	GetOne(ctx context.Context, profile, name string)
}

type emojiLogic struct {
	repo       *infra.Infra
	configRepo *infra.ConfigInfra
}

func NewEmoji(r *infra.Infra) EmojiLogic {
	return &emojiLogic{repo: r}
}

func (l emojiLogic) GetOne(ctx context.Context, profile, name string) {
	host, err := l.configRepo.GetProfileHost(profile)
	if err != nil {
		return
	}
	body := map[string]any{
		"name": name,
	}

	// https://host.tld/api/emoji
	u := (&url.URL{
		Scheme: "https",
		Host:   host,
	}).
		JoinPath("api", "emoji").String()
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(body)

	emoji, err := mi.Post2[mi.Emoji](u, buf)
	if err != nil {
		slog.Error(err.Error())
	}

	res, err := l.repo.ReactionOne(ctx, profile, name)
	if err != nil {
		// TODO: エラー処理
		slog.Error(err.Error())
	}
	l.repo.InsertEmoji(ctx, profile, res.ID, emoji)
}
