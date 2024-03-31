package logic

import (
	"context"
	"encoding/json"
	"fmt"

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
	body := map[string]any{
		"name": name,
	}
	// if id != "" {
	// 	body["untilId"] = id
	// }
	b, _ := json.Marshal(body)
	// fmt.Println(string(b))
	u := fmt.Sprintf("https://%s/api/emoji", l.repo.Config().Profiles[profile].Host)
	resp, err := mi.Post(u, b)
	if err != nil {
		fmt.Println(err)
	}
	// fmt.Println(string(resp))
	l.repo.InsertEmoji(ctx, profile, resp)
}
