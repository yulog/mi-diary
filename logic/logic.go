package logic

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/a-h/templ"
	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/app"
	cm "github.com/yulog/mi-diary/components"
	"github.com/yulog/mi-diary/model"
	"github.com/yulog/mi-diary/util/pg"
	mi "github.com/yulog/miutil"
)

type Repositorier interface {
	Reactions(ctx context.Context, profile string) ([]model.ReactionEmoji, error)
	ReactionOne(ctx context.Context, profile, name string) (model.ReactionEmoji, error)
	ReactionImageEmpty(ctx context.Context, profile string) ([]model.ReactionEmoji, error)
	HashTags(ctx context.Context, profile string) ([]model.HashTag, error)
	Users(ctx context.Context, profile string) ([]model.User, error)
	Archives(ctx context.Context, profile string) ([]model.Month, error)

	Notes(ctx context.Context, profile, s string, p *pg.Pager) ([]model.Note, error)
	ReactionNotes(ctx context.Context, profile, name string, p *pg.Pager) ([]model.Note, error)
	HashTagNotes(ctx context.Context, profile, name string, p *pg.Pager) ([]model.Note, error)
	UserNotes(ctx context.Context, profile, name string, p *pg.Pager) ([]model.Note, error)
	ArchiveNotes(ctx context.Context, profile, d string, p *pg.Pager) ([]model.Note, error)

	NoteCount(ctx context.Context, profile string) (int, error)

	UpdateEmoji(ctx context.Context, profile string, id int64, e *mi.Emoji)
	UpdateColor(ctx context.Context, profile, id, c1, c2 string)

	// TODO: bunに依存しているのは良いのか
	InsertHashTag(ctx context.Context, db bun.IDB, hashtag *model.HashTag) error
	InsertUsers(ctx context.Context, db bun.IDB, users *[]model.User) error
	InsertNotes(ctx context.Context, db bun.IDB, notes *[]model.Note) (int64, error)
	InsertReactions(ctx context.Context, db bun.IDB, reactions *[]model.ReactionEmoji) error
	InsertNoteToTags(ctx context.Context, db bun.IDB, noteToTags *[]model.NoteToTag) error
	InsertNoteToFiles(ctx context.Context, db bun.IDB, noteToFiles *[]model.NoteToFile) error
	Count(ctx context.Context, db bun.IDB) error

	RunInTx(ctx context.Context, profile string, fn func(ctx context.Context, tx bun.Tx) error)

	GenerateSchema(profile string)
	Migrate(profile string)

	// TOOD: これは良いのか
	NewFileInfra() FileRepositorier
}

type FileRepositorier interface {
	Get(ctx context.Context, profile, c string, p *pg.Pager) ([]model.File, error)
	GetByNoteID(ctx context.Context, profile, id string) ([]model.File, error)
	GetByEmptyColor(ctx context.Context, profile string) ([]model.File, error)

	Count(ctx context.Context, profile string) (int, error)

	Insert(ctx context.Context, db bun.IDB, files *[]model.File) error
}

type JobRepositorier interface {
	GetJob() chan app.Job
	SetJob(j app.Job)

	GetProgress() (int, int)
	SetProgress(p, t int) (int, int)
	UpdateProgress(p, t int) (int, int)
	GetProgressDone() bool
	SetProgressDone(d bool) bool
}

type ConfigRepositorier interface {
	SetConfig(key string, prof app.Profile)
	StoreConfig() error
	GetProfiles() *app.Profiles
	GetProfilesSorted(yield func(k string, v app.Profile) bool)
	GetProfilesSortedKey() []string
	GetProfile(key string) (app.Profile, error)
	GetProfileHost(key string) (string, error)
	GetPort() string
}

type MisskeyAPIRepositorier interface {
	GetUserReactions(profile, id string, limit int) (int, *mi.Reactions, error)
	GetEmoji(profile, name string) (*mi.Emoji, error)
}

type Logic struct {
	Repo           Repositorier
	FileRepo       FileRepositorier
	JobRepo        JobRepositorier
	ConfigRepo     ConfigRepositorier
	MisskeyAPIRepo MisskeyAPIRepositorier
}

type Dependency struct {
	repo           Repositorier
	fileRepo       FileRepositorier
	jobRepo        JobRepositorier
	configRepo     ConfigRepositorier
	misskeyAPIRepo MisskeyAPIRepositorier
}

func New() *Dependency {
	return &Dependency{}
}

func (d *Dependency) WithRepo(repo Repositorier) *Dependency {
	d.repo = repo
	return d
}

func (d *Dependency) WithFileRepo(repo FileRepositorier) *Dependency {
	d.fileRepo = repo
	return d
}

// TODO: WithFileRepoの後に使う必要がある。WithRepoはやめて、Newの引数にする？
func (d *Dependency) WithFileRepoUsingRepo() *Dependency {
	d.fileRepo = d.repo.NewFileInfra()
	return d
}

func (d *Dependency) WithJobRepo(repo JobRepositorier) *Dependency {
	d.jobRepo = repo
	return d
}

func (d *Dependency) WithConfigRepo(repo ConfigRepositorier) *Dependency {
	d.configRepo = repo
	return d
}

func (d *Dependency) WithMisskeyAPIRepo(repo MisskeyAPIRepositorier) *Dependency {
	d.misskeyAPIRepo = repo
	return d
}

func (d *Dependency) Build() *Logic {
	return &Logic{
		Repo:           d.repo,
		FileRepo:       d.fileRepo,
		JobRepo:        d.jobRepo,
		ConfigRepo:     d.configRepo,
		MisskeyAPIRepo: d.misskeyAPIRepo,
	}
}

type Params struct {
	Page  int
	S     string
	Color string
}

func (l *Logic) HomeLogic(ctx context.Context, profile string) (templ.Component, error) {
	_, err := l.ConfigRepo.GetProfile(profile)
	if err != nil {
		return nil, err
	}
	r, err := l.Repo.Reactions(ctx, profile)
	if err != nil {
		return nil, err
	}
	h, err := l.Repo.HashTags(ctx, profile)
	if err != nil {
		return nil, err
	}
	u, err := l.Repo.Users(ctx, profile)
	if err != nil {
		return nil, err
	}

	return cm.IndexParams{
		Title:     profile,
		Profile:   profile,
		Reactions: r,
		HashTags:  h,
		Users:     u,
	}.Index(), nil
}

func (l *Logic) ReactionsLogic(ctx context.Context, profile, name string, params Params) (templ.Component, error) {
	host, err := l.ConfigRepo.GetProfileHost(profile)
	if err != nil {
		return nil, err
	}

	p := pg.New(0)
	page := p.Page(params.Page)

	notes, err := l.Repo.ReactionNotes(ctx, profile, name, p)
	if err != nil {
		return nil, err
	}

	hasNext := len(notes) >= p.Limit()

	n := cm.Note{
		Title:   name,
		Profile: profile,
		Host:    host,
		Items:   notes,
	}
	cp := cm.Pages{
		Current: page,
		Prev:    cm.Page{Index: p.Prev()},
		Next:    cm.Page{Index: p.Next(), Has: hasNext},
		Last:    cm.Page{Index: p.Last()},
	}

	return n.WithPages(cp), nil
}

func (l *Logic) HashTagsLogic(ctx context.Context, profile, name string, params Params) (templ.Component, error) {
	host, err := l.ConfigRepo.GetProfileHost(profile)
	if err != nil {
		return nil, err
	}

	p := pg.New(0)
	page := p.Page(params.Page)

	notes, err := l.Repo.HashTagNotes(ctx, profile, name, p)
	if err != nil {
		return nil, err
	}

	hasNext := len(notes) >= p.Limit()

	n := cm.Note{
		Title:   name,
		Profile: profile,
		Host:    host,
		Items:   notes,
	}
	cp := cm.Pages{
		Current: page,
		Prev:    cm.Page{Index: p.Prev()},
		Next:    cm.Page{Index: p.Next(), Has: hasNext},
		Last:    cm.Page{Index: p.Last()},
	}

	return n.WithPages(cp), nil
}

func (l *Logic) UsersLogic(ctx context.Context, profile, name string, params Params) (templ.Component, error) {
	host, err := l.ConfigRepo.GetProfileHost(profile)
	if err != nil {
		return nil, err
	}

	p := pg.New(0)
	page := p.Page(params.Page)

	notes, err := l.Repo.UserNotes(ctx, profile, name, p)
	if err != nil {
		return nil, err
	}

	hasNext := len(notes) >= p.Limit()

	n := cm.Note{
		Title:   fmt.Sprintf("%s - %d", name, page),
		Profile: profile,
		Host:    host,
		Items:   notes,
	}
	cp := cm.Pages{
		Current: page,
		Prev:    cm.Page{Index: p.Prev()},
		Next:    cm.Page{Index: p.Next(), Has: hasNext},
		Last:    cm.Page{Index: p.Last()},
	}

	return n.WithPages(cp), nil
}

func (l *Logic) FilesLogic(ctx context.Context, profile string, params Params) (templ.Component, error) {
	host, err := l.ConfigRepo.GetProfileHost(profile)
	if err != nil {
		return nil, err
	}

	count := 0
	if params.Color == "" {
		count, err = l.FileRepo.Count(ctx, profile)
		if err != nil {
			return nil, err
		}
		slog.Info("File count", slog.Int("count", count))
	}

	p := pg.New(count)
	page := p.Page(params.Page)
	slog.Info("page count", slog.Int("count", page))

	files, err := l.FileRepo.Get(ctx, profile, params.Color, p)
	if err != nil {
		return nil, err
	}
	slog.Info("File result count", slog.Int("count", len(files)))
	if len(files) == 0 {
		return nil, fmt.Errorf("file not found")
	}

	hasNext := false
	if params.Color == "" {
		hasNext = len(files) >= p.Limit() && p.Next() <= p.Last()
	} else {
		hasNext = len(files) >= p.Limit()
	}
	slog.Info("has next", slog.Bool("bool", hasNext))
	hasLast := p.Next() < p.Last()
	slog.Info("has last", slog.Bool("bool", hasLast))

	n := cm.File{
		Title:          fmt.Sprint(page),
		Profile:        profile,
		Host:           host,
		FileFilterPath: fmt.Sprintf("/profiles/%s/files", profile),
		Items:          files,
	}
	cp := cm.Pages{
		Current: page,
		Prev:    cm.Page{Index: p.Prev()},
		Next:    cm.Page{Index: p.Next(), Has: hasNext},
		Last:    cm.Page{Index: p.Last(), Has: hasLast},
		QueryParams: cm.QueryParams{
			Page:  params.Page,
			Color: params.Color,
		},
	}

	return n.WithPages(cp), nil
}

func (l *Logic) NotesLogic(ctx context.Context, profile string, params Params) (templ.Component, error) {
	host, err := l.ConfigRepo.GetProfileHost(profile)
	if err != nil {
		return nil, err
	}

	count := 0
	if params.S == "" {
		count, err = l.Repo.NoteCount(ctx, profile)
		if err != nil {
			return nil, err
		}
	}

	p := pg.New(count)
	page := p.Page(params.Page)

	notes, err := l.Repo.Notes(ctx, profile, params.S, p)
	if err != nil {
		return nil, err
	}
	if len(notes) == 0 {
		return nil, fmt.Errorf("note not found")
	}
	title := ""
	if params.S != "" {
		title = fmt.Sprintf("%s - %d", params.S, page)
	} else {
		title = fmt.Sprint(page)
	}

	hasNext := false
	if params.S == "" {
		hasNext = len(notes) >= p.Limit() && p.Next() <= p.Last()
	} else {
		hasNext = len(notes) >= p.Limit()
	}
	hasLast := p.Next() < p.Last()

	n := cm.Note{
		Title:      title,
		Profile:    profile,
		Host:       host,
		SearchPath: fmt.Sprintf("/profiles/%s/notes", profile),
		Items:      notes,
	}
	cp := cm.Pages{
		Current: page,
		Prev:    cm.Page{Index: p.Prev()},
		Next:    cm.Page{Index: p.Next(), Has: hasNext},
		Last:    cm.Page{Index: p.Last(), Has: hasLast},
		QueryParams: cm.QueryParams{
			Page: params.Page,
			S:    params.S,
		},
	}

	return n.WithPages(cp), nil
}

func (l *Logic) ArchivesLogic(ctx context.Context, profile string) (templ.Component, error) {
	a, err := l.Repo.Archives(ctx, profile)
	if err != nil {
		return nil, err
	}
	return cm.ArchiveParams{
		Title:   "Archives",
		Profile: profile,
		Items:   a,
	}.Archive(), nil
}

func (l *Logic) ArchiveNotesLogic(ctx context.Context, profile, d string, params Params) (templ.Component, error) {
	host, err := l.ConfigRepo.GetProfileHost(profile)
	if err != nil {
		return nil, err
	}

	p := pg.New(0)
	page := p.Page(params.Page)

	notes, err := l.Repo.ArchiveNotes(ctx, profile, d, p)
	if err != nil {
		return nil, err
	}
	title := fmt.Sprintf("%s - %d", d, page)

	hasNext := len(notes) >= p.Limit()

	n := cm.Note{
		Title:   title,
		Profile: profile,
		Host:    host,
		Items:   notes,
	}
	cp := cm.Pages{
		Current: page,
		Prev:    cm.Page{Index: p.Prev()},
		Next:    cm.Page{Index: p.Next(), Has: hasNext},
		Last:    cm.Page{Index: p.Last()},
	}

	return n.WithPages(cp), nil
}

func (l *Logic) GenerateSchema() {
	for k := range *l.ConfigRepo.GetProfiles() {
		l.Repo.GenerateSchema(k)
		break // schemaの生成は1つだけやれば良さそう
	}
}

func (l *Logic) Migrate() {
	for k := range *l.ConfigRepo.GetProfiles() {
		l.Repo.Migrate(k)
	}
}
