package infra

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/uptrace/bun"
	"github.com/yulog/mi-diary/domain/model"
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

func (infra *Infra) Count(ctx context.Context, profile string) error {
	db, ok := txFromContext(ctx)
	if !ok {
		db = infra.DB(profile)
	}

	// 月別のカウント
	err := countMonthly(ctx, db)
	if err != nil {
		return err
	}

	// 日別のカウント
	err = countDaily(ctx, db)
	if err != nil {
		return err
	}
	return err
}

// 月別のカウント
func countMonthly(ctx context.Context, db bun.IDB) error {
	var months []model.Month
	err := db.NewSelect().
		Model((*model.Note)(nil)).
		ColumnExpr("strftime('%Y-%m', created_at, 'localtime') as ym").
		ColumnExpr("count(*) as count").
		Group("ym").
		Having("ym is not null").
		Scan(ctx, &months)
	if err != nil {
		return err
	}

	_, err = db.NewInsert().
		Model(&months).
		On("CONFLICT DO UPDATE").
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

// 日別のカウント
func countDaily(ctx context.Context, db bun.IDB) error {
	var days []model.Day
	err := db.NewSelect().
		Model((*model.Note)(nil)).
		ColumnExpr("strftime('%Y-%m-%d', created_at, 'localtime') as ymd").
		ColumnExpr("strftime('%Y-%m', created_at, 'localtime') as ym").
		ColumnExpr("count(*) as count").
		Group("ymd").
		Having("ymd is not null").
		Scan(ctx, &days)
	if err != nil {
		return err
	}

	_, err = db.NewInsert().
		Model(&days).
		On("CONFLICT DO UPDATE").
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
