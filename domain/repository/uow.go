package repository

import "context"

type UnitOfWorkRepositorier interface {
	RunInTx(ctx context.Context, profile string, fn func(ctx context.Context) error)
}
