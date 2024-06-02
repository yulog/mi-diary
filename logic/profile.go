package logic

import (
	"context"
	"fmt"
	"net"
	"net/url"

	"github.com/a-h/templ"
	"github.com/google/uuid"
	"github.com/yulog/mi-diary/app"
	cm "github.com/yulog/mi-diary/components"
	"github.com/yulog/mi-diary/migrate"
	"github.com/yulog/miutil/miauth"
)

func (l *Logic) SelectProfileLogic(ctx context.Context) templ.Component {
	var ps []string
	for k := range l.repo.Config().Profiles {
		ps = append(ps, k)
	}

	return cm.SelectProfile("Select profile...", ps)
}

func (l *Logic) NewProfileLogic(ctx context.Context) templ.Component {

	return cm.AddProfile("New Profile")
}

func (l *Logic) AddProfileLogic(ctx context.Context, server string) string {
	u, _ := url.Parse(server)

	conf := &miauth.AuthConfig{
		Name: "mi-diary-app",
		Callback: (&url.URL{
			Scheme: "http",
			Host:   net.JoinHostPort("localhost", l.repo.Config().Port),
		}).
			JoinPath("callback", u.Host).String(),
		Permission: []string{"read:reactions"},
		Host:       u.String(),
	}
	fmt.Println(conf.AuthCodeURL())

	return conf.AuthCodeURL()
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

	if resp.OK {
		cfg := l.repo.Config()
		cfg.Profiles[fmt.Sprintf("%s@%s", resp.User.Username, host)] = app.Profile{
			I:        resp.Token,
			UserId:   resp.User.ID,
			UserName: resp.User.Username,
			Host:     host,
		}
		app.ForceWriteConfig(cfg)

		migrate.Do(l.repo)
	}

	return nil
}
