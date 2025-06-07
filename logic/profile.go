package logic

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/url"

	"github.com/google/uuid"
	"github.com/yulog/mi-diary/internal/config"
	"github.com/yulog/miutil/miauth"
)

func (l *Logic) SelectProfileLogic(ctx context.Context) *SelectProfileOutput {
	return &SelectProfileOutput{
		Title:    "Select profile...",
		Profiles: l.ConfigRepo.GetProfilesSortedKey(),
	}
}

func (l *Logic) NewProfileLogic(ctx context.Context) *AddProfileOutput {

	return &AddProfileOutput{
		Title: "New Profile",
	}
}

func (l *Logic) AddProfileLogic(ctx context.Context, server string) (string, error) {
	u, err := url.Parse(server)
	if err != nil {
		return "", err
	}

	conf := &miauth.AuthConfig{
		Name: "mi-diary-app",
		Callback: (&url.URL{
			Scheme: "http",
			Host:   net.JoinHostPort("localhost", l.ConfigRepo.GetPort()),
		}).
			JoinPath("callback", u.Host).String(),
		Permission: []string{"read:reactions"},
		Host:       u.String(),
	}
	slog.Info(conf.AuthCodeURL())

	return conf.AuthCodeURL(), nil
}

func (l *Logic) CallbackLogic(ctx context.Context, host, sessionId string) error {
	id, err := uuid.Parse(sessionId)
	if err != nil {
		return err
	}

	conf := &miauth.AuthConfig{
		SessionID: id,
		Host: (&url.URL{
			Scheme: "https",
			Host:   host,
		}).String(),
	}
	resp, err := conf.Exchange(ctx)
	if err != nil {
		return err
	}

	if !resp.OK {
		return fmt.Errorf("failed to authenticate")
	}

	l.ConfigRepo.SetConfig(
		fmt.Sprintf("%s@%s", resp.User.Username, host),
		config.Profile{
			I:        resp.Token,
			UserID:   resp.User.ID,
			UserName: resp.User.Username,
			Host:     host,
		},
	)
	l.ConfigRepo.StoreConfig()

	l.Migrate()

	return nil
}
