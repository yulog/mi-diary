package logic

import (
	"context"

	"github.com/yulog/mi-diary/domain/repository"
	"github.com/yulog/mi-diary/domain/service"
	"github.com/yulog/mi-diary/util/pagination"
)

type Repositorier interface {
	// TODO: これは良いのか
	NewUnitOfWorkInfra() repository.UnitOfWorkRepositorier
	NewNoteInfra() repository.NoteRepositorier
	NewUserInfra() repository.UserRepositorier
	NewHashTagInfra() repository.HashTagRepositorier
	NewEmojiInfra() repository.EmojiRepositorier
	NewFileInfra() repository.FileRepositorier
	NewArchiveInfra() repository.ArchiveRepositorier
	NewCacheInfra() repository.CacheRepositorier
	NewMigrationInfra() service.MigrationServicer
}

type UnitOfWork interface {
	RunInTx(ctx context.Context, profile string, fn func(ctx context.Context) error)
}

type Logic struct {
	Repo             Repositorier
	UnitOfWork       UnitOfWork
	UOWRepo          repository.UnitOfWorkRepositorier
	NoteRepo         repository.NoteRepositorier
	UserRepo         repository.UserRepositorier
	HashTagRepo      repository.HashTagRepositorier
	EmojiRepo        repository.EmojiRepositorier
	FileRepo         repository.FileRepositorier
	ArchiveRepo      repository.ArchiveRepositorier
	ConfigRepo       repository.ConfigRepositorier
	CacheRepo        repository.CacheRepositorier
	MisskeyService   service.MisskeyAPIServicer
	MigrationService service.MigrationServicer
	JobWorkerService service.JobWorker
}

type Dependency struct {
	repo             Repositorier
	unitOfWork       UnitOfWork
	uowRepo          repository.UnitOfWorkRepositorier
	noteRepo         repository.NoteRepositorier
	userRepo         repository.UserRepositorier
	hashTagRepo      repository.HashTagRepositorier
	emojiRepo        repository.EmojiRepositorier
	fileRepo         repository.FileRepositorier
	archiveRepo      repository.ArchiveRepositorier
	configRepo       repository.ConfigRepositorier
	cacheRepo        repository.CacheRepositorier
	misskeyService   service.MisskeyAPIServicer
	migrationService service.MigrationServicer
	jobWorkerService service.JobWorker
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

func (d *Dependency) WithCacheRepo(repo repository.CacheRepositorier) *Dependency {
	d.cacheRepo = repo
	return d
}

func (d *Dependency) WithUOWRepoUsingRepo() *Dependency {
	d.uowRepo = d.repo.NewUnitOfWorkInfra()
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

func (d *Dependency) WithCacheRepoUsingRepo() *Dependency {
	d.cacheRepo = d.repo.NewCacheInfra()
	return d
}

func (d *Dependency) WithMigrationServiceUsingRepo() *Dependency {
	d.migrationService = d.repo.NewMigrationInfra()
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

func (d *Dependency) WithJobWorkerService(srv service.JobWorker) *Dependency {
	d.jobWorkerService = srv
	return d
}

func (d *Dependency) Build() *Logic {
	return &Logic{
		Repo:             d.repo,
		UOWRepo:          d.uowRepo,
		NoteRepo:         d.noteRepo,
		UserRepo:         d.userRepo,
		HashTagRepo:      d.hashTagRepo,
		EmojiRepo:        d.emojiRepo,
		FileRepo:         d.fileRepo,
		ArchiveRepo:      d.archiveRepo,
		ConfigRepo:       d.configRepo,
		CacheRepo:        d.cacheRepo,
		MisskeyService:   d.misskeyService,
		MigrationService: d.migrationService,
		JobWorkerService: d.jobWorkerService,
	}
}

type ItemLimitHasNextPageChecker struct {
	ItemCount int
}

func (c ItemLimitHasNextPageChecker) HasNextPage(p *pagination.Pagination) bool {
	return c.ItemCount >= p.Limit()
}
