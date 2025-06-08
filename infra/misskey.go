package infra

import (
	"github.com/yulog/mi-diary/domain/service"
	mi "github.com/yulog/miutil"
)

type MisskeyAPI struct {
	client *mi.Client
}

func NewMisskeyAPI() service.MisskeyAPIServicer {
	return &MisskeyAPI{}
}

func (infra *MisskeyAPI) Client(host, credential string) service.MisskeyAPIServicer {
	infra.client = mi.NewClient("https://"+host, credential)
	return infra
}

func (infra *MisskeyAPI) GetUserReactions(userID, untilID string, limit int) (int, *mi.Reactions, error) {
	body := map[string]any{
		"limit":  limit,
		"userId": userID,
	}
	if untilID != "" {
		body["untilId"] = untilID
	}
	req, err := infra.client.NewPostRequest("api/users/reactions", body)
	if err != nil {
		return 0, &mi.Reactions{}, err
	}
	var out mi.Reactions
	err = req.Do(&out)
	if err != nil {
		return 0, &mi.Reactions{}, err
	}
	return len(out), &out, nil
}

func (infra *MisskeyAPI) GetEmoji(name string) (*mi.Emoji, error) {
	body := map[string]any{
		"name": name,
	}
	req, err := infra.client.NewPostRequest("api/emoji", body)
	if err != nil {
		return &mi.Emoji{}, err
	}
	var out mi.Emoji
	err = req.Do(&out)
	if err != nil {
		return &mi.Emoji{}, err
	}
	return &out, nil
}
