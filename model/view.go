package model

import "github.com/uptrace/bun"

type Archive struct {
	bun.BaseModel `bun:"table:archives,alias:a"`

	YM       string
	YmCount  int64
	YMD      string
	YmdCount int64
}
