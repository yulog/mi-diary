package infra

import (
	"bytes"
	"net/url"

	"github.com/goccy/go-json"
	"github.com/yulog/mi-diary/app"
	mi "github.com/yulog/miutil"
)

func (infra *Infra) GetUserReactions(prof app.Profile, id string, limit int) (int, *mi.Reactions, error) {
	body := map[string]any{
		"i":      prof.I,
		"limit":  limit,
		"userId": prof.UserId,
	}
	if id != "" {
		body["untilId"] = id
	}

	// https://host.tld/api/users/reactions
	// 却って分かりにくい気もする
	u := (&url.URL{
		Scheme: "https",
		Host:   prof.Host,
	}).
		JoinPath("api", "users", "reactions").String()

	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(body)

	r, err := mi.Post2[mi.Reactions](u, buf)
	if err != nil {
		return 0, &mi.Reactions{}, err
	}

	return len(*r), r, nil
}

func (infra *Infra) GetEmoji(prof app.Profile, name string) (*mi.Emoji, error) {
	body := map[string]any{
		"name": name,
	}

	// https://host.tld/api/emoji
	u := (&url.URL{
		Scheme: "https",
		Host:   prof.Host,
	}).
		JoinPath("api", "emoji").String()
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(body)

	emoji, err := mi.Post2[mi.Emoji](u, buf)
	if err != nil {
		return &mi.Emoji{}, err
	}

	return emoji, nil
}
