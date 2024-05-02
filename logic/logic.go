package logic

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"time"

	"github.com/a-h/templ"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/yulog/mi-diary/app"
	cm "github.com/yulog/mi-diary/components"
	"github.com/yulog/mi-diary/infra"
	"github.com/yulog/mi-diary/migrate"
	"github.com/yulog/mi-diary/util/pg"
	mi "github.com/yulog/miutil"
	"github.com/yulog/miutil/miauth"
)

type Logic struct {
	repo *infra.Infra
}

func New(r *infra.Infra) *Logic {
	return &Logic{repo: r}
}

func (l *Logic) SelectProfileLogic(ctx context.Context) templ.Component {
	var ps []string
	for k := range l.repo.Config().Profiles {
		ps = append(ps, k)
	}

	return cm.SelectProfile("Select profile...", ps)
}

func (l *Logic) HomeLogic(ctx context.Context, profile string) (templ.Component, error) {

	if _, ok := l.repo.Config().Profiles[profile]; !ok {
		return nil, fmt.Errorf("invalid profile: %s", profile)
	}

	return cm.IndexParams{
		Title:     "Home",
		Profile:   profile,
		Reactions: l.repo.Reactions(ctx, profile),
		HashTags:  l.repo.HashTags(ctx, profile),
		Users:     l.repo.Users(ctx, profile),
	}.Index(), nil
}

func (l *Logic) ReactionsLogic(ctx context.Context, profile, name string, page int) templ.Component {
	p := pg.New(0)
	page = p.Page(page)

	notes := l.repo.ReactionNotes(ctx, profile, name, p)

	hasNext := len(notes) >= p.Limit()

	n := cm.Note{
		Title:   name,
		Profile: profile,
		Host:    l.repo.Config().Profiles[profile].Host,
		Items:   notes,
	}
	cp := cm.Pages{
		Current: page,
		Prev:    p.Prev(),
		Next:    p.Next(),
		Last:    p.Last(),
		HasNext: hasNext,
		HasLast: false,
	}

	return n.WithPages(cp)
}

func (l *Logic) HashTagsLogic(ctx context.Context, profile, name string, page int) templ.Component {
	p := pg.New(0)
	page = p.Page(page)

	notes := l.repo.HashTagNotes(ctx, profile, name, p)

	hasNext := len(notes) >= p.Limit()

	n := cm.Note{
		Title:   name,
		Profile: profile,
		Host:    l.repo.Config().Profiles[profile].Host,
		Items:   notes,
	}
	cp := cm.Pages{
		Current: page,
		Prev:    p.Prev(),
		Next:    p.Next(),
		Last:    p.Last(),
		HasNext: hasNext,
		HasLast: false,
	}

	return n.WithPages(cp)
}

func (l *Logic) UsersLogic(ctx context.Context, profile, name string, page int) templ.Component {
	p := pg.New(0)
	page = p.Page(page)

	notes := l.repo.UserNotes(ctx, profile, name, p)

	hasNext := len(notes) >= p.Limit()

	n := cm.Note{
		Title:   fmt.Sprintf("%s - %d", name, page),
		Profile: profile,
		Host:    l.repo.Config().Profiles[profile].Host,
		Items:   notes,
	}
	cp := cm.Pages{
		Current: page,
		Prev:    p.Prev(),
		Next:    p.Next(),
		Last:    p.Last(),
		HasNext: hasNext,
		HasLast: false,
	}

	return n.WithPages(cp)
}

func (l *Logic) FilesLogic(ctx context.Context, profile string, page int) (templ.Component, error) {
	count, err := l.repo.FileCount(ctx, profile)
	if err != nil {
		return nil, err
	}
	fmt.Println(count)
	p := pg.New(count)
	page = p.Page(page)

	files := l.repo.Files(ctx, profile, p)
	if len(files) == 0 {
		return nil, fmt.Errorf("file not found")
	}

	hasNext := len(files) >= p.Limit() && p.Next() <= p.Last()
	hasLast := p.Next() < p.Last()

	n := cm.File{
		Title:   fmt.Sprint(page),
		Profile: profile,
		Host:    l.repo.Config().Profiles[profile].Host,
		Items:   files,
	}
	cp := cm.Pages{
		Current: page,
		Prev:    p.Prev(),
		Next:    p.Next(),
		Last:    p.Last(),
		HasNext: hasNext,
		HasLast: hasLast,
	}

	return n.WithPages(cp), nil
}

func (l *Logic) NotesLogic(ctx context.Context, profile string, page int) (templ.Component, error) {
	count, err := l.repo.NoteCount(ctx, profile)
	if err != nil {
		return nil, err
	}
	p := pg.New(count)
	page = p.Page(page)

	notes := l.repo.Notes(ctx, profile, p)
	if len(notes) == 0 {
		return nil, fmt.Errorf("note not found")
	}
	// title := fmt.Sprint(page)

	hasNext := len(notes) >= p.Limit() && p.Next() <= p.Last()
	hasLast := p.Next() < p.Last()

	n := cm.Note{
		Title:   fmt.Sprint(page),
		Profile: profile,
		Host:    l.repo.Config().Profiles[profile].Host,
		Items:   notes,
	}
	cp := cm.Pages{
		Current: page,
		Prev:    p.Prev(),
		Next:    p.Next(),
		Last:    p.Last(),
		HasNext: hasNext,
		HasLast: hasLast,
	}

	return n.WithPages(cp), nil
}

func (l *Logic) ArchivesLogic(ctx context.Context, profile string) templ.Component {
	return cm.ArchiveParams{
		Title:   "Archives",
		Profile: profile,
		Items:   l.repo.Archives(ctx, profile),
	}.Archive()
}

var reym = regexp.MustCompile(`^\d{4}-\d{2}$`)
var reymd = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func (l *Logic) ArchiveNotesLogic(ctx context.Context, profile, d string, page int) templ.Component {
	col := ""
	if reym.MatchString(d) {
		col = "strftime('%Y-%m', created_at, 'localtime')"
	} else if reymd.MatchString(d) {
		col = "strftime('%Y-%m-%d', created_at, 'localtime')"
	}

	p := pg.New(0)
	page = p.Page(page)

	notes := l.repo.ArchiveNotes(ctx, profile, col, d, p)
	title := fmt.Sprintf("%s - %d", d, page)

	hasNext := len(notes) >= p.Limit()

	n := cm.Note{
		Title:   title,
		Profile: profile,
		Host:    l.repo.Config().Profiles[profile].Host,
		Items:   notes,
	}
	cp := cm.Pages{
		Current: page,
		Prev:    p.Prev(),
		Next:    p.Next(),
		Last:    p.Last(),
		HasNext: hasNext,
		HasLast: false,
	}

	return n.WithPages(cp)
}

func (l *Logic) ManageLogic(ctx context.Context) templ.Component {
	p := l.repo.GetProgress()
	var ps []string
	for k := range l.repo.Config().Profiles {
		ps = append(ps, k)
	}
	if p > 0 {
		return cm.ManageStart("Manage")
	}
	return cm.ManageInit("Manage", ps)
}

func (l *Logic) JobStartLogic(ctx context.Context, job app.Job) templ.Component {
	l.repo.SetJob(job)

	return cm.Start("/job", fmt.Sprintf("/profiles/%s/settings/reactions", "profile"), "Reaction ID", "reaction-id", "", "Get", job.Profile, job.Type.String(), job.ID)
}

func (l *Logic) JobProgressLogic(ctx context.Context) (int, templ.Component) {
	p := l.repo.GetProgress()
	return p, cm.Progress(p)
}

func (l *Logic) JobLogic(ctx context.Context, profile string) templ.Component {
	res := l.repo.GetProgress()
	l.repo.SetProgress(0)
	var ps []string
	for k := range l.repo.Config().Profiles {
		ps = append(ps, k)
	}

	return cm.Job("/job", "/job/progress", fmt.Sprintf("/profiles/%s/settings/reactions", profile), "Reaction ID", "reaction-id", "", "Get", res, ps)
}

func (l *Logic) JobProcesser() {
	// for j := range job {
	for j := range l.repo.GetJob() {
		for i := 0; i < 10; i++ {
			p := l.repo.GetProgress()
			p = l.repo.SetProgress(p + 10)
			fmt.Println(j, p)
			time.Sleep(time.Second)
		}
	}
}

func (l *Logic) SettingsLogic(ctx context.Context, profile string) templ.Component {

	return cm.Settings("settings", profile)
}

func (l *Logic) SettingsReactionsLogic(ctx context.Context, profile, id string) {
	body := map[string]any{
		"i":      l.repo.Config().Profiles[profile].I,
		"limit":  20,
		"userId": l.repo.Config().Profiles[profile].UserId,
	}
	if id != "" {
		body["untilId"] = id
	}

	// https://host.tld/api/users/reactions
	// 却って分かりにくい気もする
	u := (&url.URL{
		Scheme: "https",
		Host:   l.repo.Config().Profiles[profile].Host,
	}).
		JoinPath("api", "users", "reactions").String()

	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(body)

	r, err := mi.Post2[mi.Reactions](u, buf)
	if err != nil {
		fmt.Println(err)
	}

	l.repo.Insert(ctx, profile, r)
}

func (l *Logic) SettingsEmojisLogic(ctx context.Context, profile, name string) {
	body := map[string]any{
		"name": name,
	}

	// https://host.tld/api/emoji
	u := (&url.URL{
		Scheme: "https",
		Host:   l.repo.Config().Profiles[profile].Host,
	}).
		JoinPath("api", "emoji").String()
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(body)

	emoji, err := mi.Post2[mi.Emoji](u, buf)
	if err != nil {
		fmt.Println(err)
	}

	l.repo.InsertEmoji(ctx, profile, emoji)
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
