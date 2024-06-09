package logic

import (
	"bytes"
	"context"
	"fmt"
	"net/url"

	"github.com/goccy/go-json"
	"github.com/yulog/mi-diary/infra"
	mi "github.com/yulog/miutil"
)

type EmojiLogic interface {
	GetOne(ctx context.Context, profile, name string)
}

type emojiLogic struct {
	repo *infra.Infra
}

func NewEmoji(r *infra.Infra) EmojiLogic {
	return &emojiLogic{repo: r}
}

func (l emojiLogic) GetOne(ctx context.Context, profile, name string) {
	host, err := l.repo.GetProfileHost(profile)
	if err != nil {
		return
	}
	body := map[string]any{
		"name": name,
	}

	// b, _ := json.Marshal(body)
	// fmt.Println(string(b))
	// https://host.tld/api/emoji
	// u := fmt.Sprintf("https://%s/api/emoji", l.repo.Config().Profiles[profile].Host)
	u := (&url.URL{
		Scheme: "https",
		Host:   host,
	}).
		JoinPath("api", "emoji").String()
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(body)

	emoji, err := mi.Post2[mi.Emoji](u, buf)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(string(resp))
	res, err := l.repo.ReactionOne(ctx, profile, name)
	if err != nil {
		// TODO: エラー処理
		fmt.Println(err)
	}
	l.repo.InsertEmoji(ctx, profile, res.ID, emoji)
}
