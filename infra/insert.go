package infra

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/uptrace/bun"
)

// func (infra *Infra) InsertFromFile(ctx context.Context, profile string) {
// 	// JSON読み込み
// 	f, _ := os.ReadFile("users_reactions.json")
// 	var r mi.Reactions
// 	json.Unmarshal(f, &r)

// 	infra.Insert(ctx, profile, &r)
// }

type txKey struct{}

func txFromContext(ctx context.Context) (bun.IDB, bool) {
	tx, ok := ctx.Value(txKey{}).(bun.IDB)
	return tx, ok
}

func (infra *Infra) RunInTx(ctx context.Context, profile string, fn func(ctx context.Context) error) {
	err := infra.DB(profile).RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		ctx = context.WithValue(ctx, txKey{}, tx)
		if err := fn(ctx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}
