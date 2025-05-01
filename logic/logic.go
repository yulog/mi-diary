package logic

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/a-h/templ"
	cm "github.com/yulog/mi-diary/components"
	"github.com/yulog/mi-diary/domain/repository"
	"github.com/yulog/mi-diary/domain/service"
	"github.com/yulog/mi-diary/util/pagination"
)

type Repositorier interface {
	RunInTx(ctx context.Context, profile string, fn func(ctx context.Context) error)

	GenerateSchema(profile string)
	Migrate(profile string)

	// TODO: これは良いのか
	NewNoteInfra() repository.NoteRepositorier
	NewUserInfra() repository.UserRepositorier
	NewHashTagInfra() repository.HashTagRepositorier
	NewEmojiInfra() repository.EmojiRepositorier
	NewFileInfra() repository.FileRepositorier
	NewArchiveInfra() repository.ArchiveRepositorier
}

type Logic struct {
	Repo           Repositorier
	NoteRepo       repository.NoteRepositorier
	UserRepo       repository.UserRepositorier
	HashTagRepo    repository.HashTagRepositorier
	EmojiRepo      repository.EmojiRepositorier
	FileRepo       repository.FileRepositorier
	ArchiveRepo    repository.ArchiveRepositorier
	JobRepo        repository.JobRepositorier
	ConfigRepo     repository.ConfigRepositorier
	MisskeyService service.MisskeyAPIServicer
}

type Dependency struct {
	repo           Repositorier
	noteRepo       repository.NoteRepositorier
	userRepo       repository.UserRepositorier
	hashTagRepo    repository.HashTagRepositorier
	emojiRepo      repository.EmojiRepositorier
	fileRepo       repository.FileRepositorier
	archiveRepo    repository.ArchiveRepositorier
	jobRepo        repository.JobRepositorier
	configRepo     repository.ConfigRepositorier
	misskeyService service.MisskeyAPIServicer
}

func New() *Dependency {
	return &Dependency{}
}

func (d *Dependency) WithRepo(repo Repositorier) *Dependency {
	d.repo = repo
	return d
}

func (d *Dependency) WithNoteRepo(repo repository.NoteRepositorier) *Dependency {
	d.noteRepo = repo
	return d
}

func (d *Dependency) WithUserRepo(repo repository.UserRepositorier) *Dependency {
	d.userRepo = repo
	return d
}

func (d *Dependency) WithHashTagRepo(repo repository.HashTagRepositorier) *Dependency {
	d.hashTagRepo = repo
	return d
}

func (d *Dependency) WithEmojiRepo(repo repository.EmojiRepositorier) *Dependency {
	d.emojiRepo = repo
	return d
}

func (d *Dependency) WithFileRepo(repo repository.FileRepositorier) *Dependency {
	d.fileRepo = repo
	return d
}

func (d *Dependency) WithArchiveRepo(repo repository.ArchiveRepositorier) *Dependency {
	d.archiveRepo = repo
	return d
}

func (d *Dependency) WithNoteRepoUsingRepo() *Dependency {
	d.noteRepo = d.repo.NewNoteInfra()
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

func (d *Dependency) WithArchiveRepoUsingRepo() *Dependency {
	d.archiveRepo = d.repo.NewArchiveInfra()
	return d
}

func (d *Dependency) WithJobRepo(repo repository.JobRepositorier) *Dependency {
	d.jobRepo = repo
	return d
}

func (d *Dependency) WithConfigRepo(repo repository.ConfigRepositorier) *Dependency {
	d.configRepo = repo
	return d
}

func (d *Dependency) WithMisskeyAPIRepo(srv service.MisskeyAPIServicer) *Dependency {
	d.misskeyService = srv
	return d
}

func (d *Dependency) Build() *Logic {
	return &Logic{
		Repo:           d.repo,
		NoteRepo:       d.noteRepo,
		UserRepo:       d.userRepo,
		HashTagRepo:    d.hashTagRepo,
		EmojiRepo:      d.emojiRepo,
		FileRepo:       d.fileRepo,
		ArchiveRepo:    d.archiveRepo,
		JobRepo:        d.jobRepo,
		ConfigRepo:     d.configRepo,
		MisskeyService: d.misskeyService,
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

	notes, err := l.NoteRepo.GetByReaction(ctx, profile, name, p)
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

	notes, err := l.NoteRepo.GetByHashTag(ctx, profile, name, p)
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

	notes, err := l.NoteRepo.GetByUser(ctx, profile, name, p)
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
		count, err = l.NoteRepo.Count(ctx, profile)
		if err != nil {
			return nil, err
		}
	}

	p, err := pagination.New(params.Page, 10, count, nil, nil)
	if err != nil {
		slog.Error(err.Error())
	}

	notes, err := l.NoteRepo.Get(ctx, profile, params.S, p)
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
	a, err := l.ArchiveRepo.Get(ctx, profile)
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

	notes, err := l.NoteRepo.GetByArchive(ctx, profile, d, p)
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
