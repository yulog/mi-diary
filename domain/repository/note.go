package repository

import (
	"context"

	"github.com/yulog/mi-diary/domain/model"
	"github.com/yulog/mi-diary/util/pagination"
)

type NoteRepositorier interface {
	Get(ctx context.Context, profile, s string, p pagination.Paging) ([]model.Note, error)
	GetByReaction(ctx context.Context, profile, name string, p pagination.Paging) ([]model.Note, error)
	GetByHashTag(ctx context.Context, profile, name string, p pagination.Paging) ([]model.Note, error)
	GetByUser(ctx context.Context, profile, name string, p pagination.Paging) ([]model.Note, error)
	GetByArchive(ctx context.Context, profile, d string, p pagination.Paging) ([]model.Note, error)

	Count(ctx context.Context, profile string) (int, error)

	Insert(ctx context.Context, profile string, notes *[]model.Note) (int64, error)
	InsertNoteToTags(ctx context.Context, profile string, noteToTags *[]model.NoteToTag) error
	InsertNoteToFiles(ctx context.Context, profile string, noteToFiles *[]model.NoteToFile) error
}
