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
	"github.com/yulog/mi-diary/util/pagination"
	mi "github.com/yulog/miutil"
)

type Repositorier interface {
	Archives(ctx context.Context, profile string) ([]model.Month, error)

	Notes(ctx context.Context, profile, s string, p pagination.Paging) ([]model.Note, error)
	ReactionNotes(ctx context.Context, profile, name string, p pagination.Paging) ([]model.Note, error)
	HashTagNotes(ctx context.Context, profile, name string, p pagination.Paging) ([]model.Note, error)
	UserNotes(ctx context.Context, profile, name string, p pagination.Paging) ([]model.Note, error)
	ArchiveNotes(ctx context.Context, profile, d string, p pagination.Paging) ([]model.Note, error)

	NoteCount(ctx context.Context, profile string) (int, error)

	// TODO: bunに依存しているのは良いのか
	InsertNotes(ctx context.Context, db bun.IDB, notes *[]model.Note) (int64, error)
	InsertNoteToTags(ctx context.Context, db bun.IDB, noteToTags *[]model.NoteToTag) error
	InsertNoteToFiles(ctx context.Context, db bun.IDB, noteToFiles *[]model.NoteToFile) error
	Count(ctx context.Context, db bun.IDB) error

	RunInTx(ctx context.Context, profile string, fn func(ctx context.Context, tx bun.Tx) error)

	GenerateSchema(profile string)
	Migrate(profile string)

	// TODO: これは良いのか
	NewUserInfra() UserRepositorier
	NewHashTagInfra() HashTagRepositorier
	NewEmojiInfra() EmojiRepositorier
	NewFileInfra() FileRepositorier
}

type HashTagRepositorier interface {
	Get(ctx context.Context, profile string) ([]model.HashTag, error)

	Insert(ctx context.Context, db bun.IDB, hashtag *model.HashTag) error
}

type EmojiRepositorier interface {
	Get(ctx context.Context, profile string) ([]model.ReactionEmoji, error)
	GetByName(ctx context.Context, profile, name string) (model.ReactionEmoji, error)
	GetByEmptyImage(ctx context.Context, profile string) ([]model.ReactionEmoji, error)

	Insert(ctx context.Context, db bun.IDB, reactions *[]model.ReactionEmoji) error

	UpdateByPKWithImage(ctx context.Context, profile string, id int64, e *mi.Emoji)
}

type FileRepositorier interface {
	Get(ctx context.Context, profile, c string, p pagination.Paging) ([]model.File, error)
	GetByNoteID(ctx context.Context, profile, id string) ([]model.File, error)
	GetByEmptyColor(ctx context.Context, profile string) ([]model.File, error)

	Count(ctx context.Context, profile string) (int, error)

	Insert(ctx context.Context, db bun.IDB, files *[]model.File) error

	UpdateByPKWithColor(ctx context.Context, profile, id, c1, c2 string)
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
	UserRepo       UserRepositorier
	HashTagRepo    HashTagRepositorier
	EmojiRepo      EmojiRepositorier
	FileRepo       FileRepositorier
	JobRepo        JobRepositorier
	ConfigRepo     ConfigRepositorier
	MisskeyAPIRepo MisskeyAPIRepositorier
}

type Dependency struct {
	repo           Repositorier
	userRepo       UserRepositorier
	hashTagRepo    HashTagRepositorier
	emojiRepo      EmojiRepositorier
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

func (d *Dependency) WithUserRepo(repo UserRepositorier) *Dependency {
	d.userRepo = repo
	return d
}

func (d *Dependency) WithHashTagRepo(repo HashTagRepositorier) *Dependency {
	d.hashTagRepo = repo
	return d
}

func (d *Dependency) WithEmojiRepo(repo EmojiRepositorier) *Dependency {
	d.emojiRepo = repo
	return d
}

func (d *Dependency) WithFileRepo(repo FileRepositorier) *Dependency {
	d.fileRepo = repo
	return d
}

func (d *Dependency) WithUserRepoUsingRepo() *Dependency {
	d.userRepo = d.repo.NewUserInfra()
	return d
}

func (d *Dependency) WithHashTagRepoUsingRepo() *Dependency {
	d.hashTagRepo = d.repo.NewHashTagInfra()
	return d
}

func (d *Dependency) WithEmojiRepoUsingRepo() *Dependency {
	d.emojiRepo = d.repo.NewEmojiInfra()
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
		UserRepo:       d.userRepo,
		HashTagRepo:    d.hashTagRepo,
		EmojiRepo:      d.emojiRepo,
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

type ItemLimitHasNextPageChecker struct {
	ItemCount int
}

func (c ItemLimitHasNextPageChecker) HasNextPage(p *pagination.Pagination) bool {
	return c.ItemCount >= p.Limit()
}

func (l *Logic) HomeLogic(ctx context.Context, profile string) (templ.Component, error) {
	_, err := l.ConfigRepo.GetProfile(profile)
	if err != nil {
		return nil, err
	}
	r, err := l.EmojiRepo.Get(ctx, profile)
	if err != nil {
		return nil, err
	}

	return cm.IndexParams{
		Title:     profile,
		Profile:   profile,
		Reactions: r,
	}.Index(), nil
}

func (l *Logic) HashTagsLogic(ctx context.Context, profile string) (templ.Component, error) {
	_, err := l.ConfigRepo.GetProfile(profile)
	if err != nil {
		return nil, err
	}
	h, err := l.HashTagRepo.Get(ctx, profile)
	if err != nil {
		return nil, err
	}
	return cm.HashTags(profile, h), nil
}

func (l *Logic) ReactionNotesLogic(ctx context.Context, profile, name string, params Params) (templ.Component, error) {
	host, err := l.ConfigRepo.GetProfileHost(profile)
	if err != nil {
		return nil, err
	}

	p, err := pagination.New(params.Page, 10, 0, nil, nil)
	if err != nil {
		slog.Error(err.Error())
	}

	notes, err := l.Repo.ReactionNotes(ctx, profile, name, p)
	if err != nil {
		return nil, err
	}

	p.NextChecker = ItemLimitHasNextPageChecker{ItemCount: len(notes)}

	next, err := p.NextPage()
	if err != nil {
		slog.Info(err.Error())
	}
	prev, err := p.PreviousPage()
	if err != nil {
		slog.Info(err.Error())
	}

	n := cm.Note{
		Title:   name,
		Profile: profile,
		Host:    host,
		Items:   notes,
	}
	cp := cm.Pages{
		Current: p.CurrentPage,
		Prev:    cm.Page{Index: prev, Has: p.HasPreviousPage()},
		Next:    cm.Page{Index: next, Has: p.HasNextPage()},
		Last:    cm.Page{Index: p.TotalPages()},
	}

	return n.WithPages(cp), nil
}

func (l *Logic) HashTagNotesLogic(ctx context.Context, profile, name string, params Params) (templ.Component, error) {
	host, err := l.ConfigRepo.GetProfileHost(profile)
	if err != nil {
		return nil, err
	}

	p, err := pagination.New(params.Page, 10, 0, nil, nil)
	if err != nil {
		slog.Error(err.Error())
	}

	notes, err := l.Repo.HashTagNotes(ctx, profile, name, p)
	if err != nil {
		return nil, err
	}

	p.NextChecker = ItemLimitHasNextPageChecker{ItemCount: len(notes)}

	next, err := p.NextPage()
	if err != nil {
		slog.Info(err.Error())
	}
	prev, err := p.PreviousPage()
	if err != nil {
		slog.Info(err.Error())
	}

	n := cm.Note{
		Title:   name,
		Profile: profile,
		Host:    host,
		Items:   notes,
	}
	cp := cm.Pages{
		Current: p.CurrentPage,
		Prev:    cm.Page{Index: prev, Has: p.HasPreviousPage()},
		Next:    cm.Page{Index: next, Has: p.HasNextPage()},
		Last:    cm.Page{Index: p.TotalPages()},
	}

	return n.WithPages(cp), nil
}

func (l *Logic) UserLogic(ctx context.Context, profile, name string, params Params) (templ.Component, error) {
	host, err := l.ConfigRepo.GetProfileHost(profile)
	if err != nil {
		return nil, err
	}

	p, err := pagination.New(params.Page, 10, 0, nil, nil)
	if err != nil {
		slog.Error(err.Error())
	}

	notes, err := l.Repo.UserNotes(ctx, profile, name, p)
	if err != nil {
		return nil, err
	}

	p.NextChecker = ItemLimitHasNextPageChecker{ItemCount: len(notes)}

	next, err := p.NextPage()
	if err != nil {
		slog.Info(err.Error())
	}
	prev, err := p.PreviousPage()
	if err != nil {
		slog.Info(err.Error())
	}

	n := cm.Note{
		Title:   fmt.Sprintf("%s - %d", name, p.CurrentPage),
		Profile: profile,
		Host:    host,
		Items:   notes,
	}
	cp := cm.Pages{
		Current: p.CurrentPage,
		Prev:    cm.Page{Index: prev, Has: p.HasPreviousPage()},
		Next:    cm.Page{Index: next, Has: p.HasNextPage()},
		Last:    cm.Page{Index: p.TotalPages()},
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

	p, err := pagination.New(params.Page, 10, count, nil, nil)
	if err != nil {
		slog.Error(err.Error())
	}
	slog.Info("page count", slog.Int("count", p.CurrentPage))

	files, err := l.FileRepo.Get(ctx, profile, params.Color, p)
	if err != nil {
		return nil, err
	}
	slog.Info("file result count", slog.Int("count", len(files)))
	if len(files) == 0 {
		return nil, fmt.Errorf("file not found")
	}

	if params.Color != "" {
		p.NextChecker = ItemLimitHasNextPageChecker{ItemCount: len(files)}
	}
	slog.Info("has next", slog.Bool("bool", p.HasNextPage()))

	next, err := p.NextPage()
	if err != nil {
		slog.Info(err.Error())
	}
	prev, err := p.PreviousPage()
	if err != nil {
		slog.Info(err.Error())
	}

	hasLast := p.CurrentPage+1 < p.TotalPages()
	slog.Info("has last", slog.Bool("bool", hasLast))

	n := cm.File{
		Title:          fmt.Sprint(p.CurrentPage),
		Profile:        profile,
		Host:           host,
		FileFilterPath: fmt.Sprintf("/profiles/%s/files", profile),
		Items:          files,
	}
	cp := cm.Pages{
		Current: p.CurrentPage,
		Prev:    cm.Page{Index: prev, Has: p.HasPreviousPage()},
		Next:    cm.Page{Index: next, Has: p.HasNextPage()},
		Last:    cm.Page{Index: p.TotalPages(), Has: hasLast},
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

	p, err := pagination.New(params.Page, 10, count, nil, nil)
	if err != nil {
		slog.Error(err.Error())
	}

	notes, err := l.Repo.Notes(ctx, profile, params.S, p)
	if err != nil {
		return nil, err
	}
	if len(notes) == 0 {
		return nil, fmt.Errorf("note not found")
	}
	title := ""
	if params.S != "" {
		title = fmt.Sprintf("%s - %d", params.S, p.CurrentPage)
		p.NextChecker = ItemLimitHasNextPageChecker{ItemCount: len(notes)}
	} else {
		title = fmt.Sprint(p.CurrentPage)
	}

	next, err := p.NextPage()
	if err != nil {
		slog.Info(err.Error())
	}
	prev, err := p.PreviousPage()
	if err != nil {
		slog.Info(err.Error())
	}

	hasLast := p.CurrentPage+1 < p.TotalPages()
	// slog.Info("has last", slog.Bool("hasLast", hasLast), slog.Int("next", next), slog.Int("total", p.TotalPages()), slog.Int("current", p.CurrentPage))

	n := cm.Note{
		Title:      title,
		Profile:    profile,
		Host:       host,
		SearchPath: fmt.Sprintf("/profiles/%s/notes", profile),
		Items:      notes,
	}
	cp := cm.Pages{
		Current: p.CurrentPage,
		Prev:    cm.Page{Index: prev, Has: p.HasPreviousPage()},
		Next:    cm.Page{Index: next, Has: p.HasNextPage()},
		Last:    cm.Page{Index: p.TotalPages(), Has: hasLast},
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

	p, err := pagination.New(params.Page, 10, 0, nil, nil)
	if err != nil {
		slog.Error(err.Error())
	}
	// slog.Info("perPage", slog.Int("perPage", p2.Limit()))

	notes, err := l.Repo.ArchiveNotes(ctx, profile, d, p)
	if err != nil {
		return nil, err
	}

	p.NextChecker = ItemLimitHasNextPageChecker{ItemCount: len(notes)}

	next, err := p.NextPage()
	if err != nil {
		slog.Info(err.Error())
	}
	prev, err := p.PreviousPage()
	if err != nil {
		slog.Info(err.Error())
	}

	n := cm.Note{
		Title:   fmt.Sprintf("%s - %d", d, p.CurrentPage),
		Profile: profile,
		Host:    host,
		Items:   notes,
	}
	cp := cm.Pages{
		Current: p.CurrentPage,
		Prev:    cm.Page{Index: prev, Has: p.HasPreviousPage()},
		Next:    cm.Page{Index: next, Has: p.HasNextPage()},
		Last:    cm.Page{Index: p.TotalPages()},
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
