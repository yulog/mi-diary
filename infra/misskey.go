package infra

import (
	"bytes"
	"net/url"

	"github.com/goccy/go-json"
	"github.com/yulog/mi-diary/app"
	"github.com/yulog/mi-diary/domain/service"
	mi "github.com/yulog/miutil"
)

type MisskeyAPI struct {
	app *app.App
}

func NewMisskeyAPI(a *app.App) service.MisskeyAPIServicer {
	return &MisskeyAPI{app: a}
}

func (infra *MisskeyAPI) GetUserReactions(profile, id string, limit int) (int, *mi.Reactions, error) {
	prof, err := infra.app.Config.Profiles.Get(profile)
	if err != nil {
		return 0, &mi.Reactions{}, err
	}
	body := map[string]any{
		"i":      prof.I,
		"limit":  limit,
		"userId": prof.UserID,
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

func (infra *MisskeyAPI) GetEmoji(profile, name string) (*mi.Emoji, error) {
	prof, err := infra.app.Config.Profiles.Get(profile)
	if err != nil {
		return &mi.Emoji{}, err
	}
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
