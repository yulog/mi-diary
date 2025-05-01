package logic

import (
	"context"

	"github.com/yulog/mi-diary/domain/repository"
	"github.com/yulog/mi-diary/domain/service"
	"github.com/yulog/mi-diary/util/pagination"
)

type Repositorier interface {
	RunInTx(ctx context.Context, profile string, fn func(ctx context.Context) error)

	// TODO: これは良いのか
	NewNoteInfra() repository.NoteRepositorier
	NewUserInfra() repository.UserRepositorier
	NewHashTagInfra() repository.HashTagRepositorier
	NewEmojiInfra() repository.EmojiRepositorier
	NewFileInfra() repository.FileRepositorier
	NewArchiveInfra() repository.ArchiveRepositorier
	NewMigrationInfra() service.MigrationServicer
}

type Logic struct {
	Repo             Repositorier
	NoteRepo         repository.NoteRepositorier
	UserRepo         repository.UserRepositorier
	HashTagRepo      repository.HashTagRepositorier
	EmojiRepo        repository.EmojiRepositorier
	FileRepo         repository.FileRepositorier
	ArchiveRepo      repository.ArchiveRepositorier
	JobRepo          repository.JobRepositorier
	ConfigRepo       repository.ConfigRepositorier
	MisskeyService   service.MisskeyAPIServicer
	MigrationService service.MigrationServicer
}

type Dependency struct {
	repo             Repositorier
	noteRepo         repository.NoteRepositorier
	userRepo         repository.UserRepositorier
	hashTagRepo      repository.HashTagRepositorier
	emojiRepo        repository.EmojiRepositorier
	fileRepo         repository.FileRepositorier
	archiveRepo      repository.ArchiveRepositorier
	jobRepo          repository.JobRepositorier
	configRepo       repository.ConfigRepositorier
	misskeyService   service.MisskeyAPIServicer
	migrationService service.MigrationServicer
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

func (d *Dependency) WithFileRepoUsingRepo() *Dependency {
	d.fileRepo = d.repo.NewFileInfra()
	return d
}

func (d *Dependency) WithArchiveRepoUsingRepo() *Dependency {
	d.archiveRepo = d.repo.NewArchiveInfra()
	return d
}

func (d *Dependency) WithMigrationServiceUsingRepo() *Dependency {
	d.migrationService = d.repo.NewMigrationInfra()
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

func (d *Dependency) WithMigrationService(srv service.MigrationServicer) *Dependency {
	d.migrationService = srv
	return d
}

func (d *Dependency) Build() *Logic {
	return &Logic{
		Repo:             d.repo,
		NoteRepo:         d.noteRepo,
		UserRepo:         d.userRepo,
		HashTagRepo:      d.hashTagRepo,
		EmojiRepo:        d.emojiRepo,
		FileRepo:         d.fileRepo,
		ArchiveRepo:      d.archiveRepo,
		JobRepo:          d.jobRepo,
		ConfigRepo:       d.configRepo,
		MisskeyService:   d.misskeyService,
		MigrationService: d.migrationService,
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
