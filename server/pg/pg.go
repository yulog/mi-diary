package pg

import (
	"math"
	"sync"

	"github.com/labstack/echo/v4"
)

type Pager struct {
	c    *echo.Context
	page int
	// Prev int
	// Next int
	lastOnce sync.Once
	last     int

	total int

	limit  int
	offset int
}

func New(c *echo.Context, total int) *Pager {
	return &Pager{
		c: c,

		total: total,
		limit: 10,
	}
}

func (p *Pager) Page() (int, error) {
	p.page = 1
	if err := echo.QueryParamsBinder(*p.c).
		Int("page", &p.page).
		BindError(); err != nil {
		return p.page, err
	}
	if p.page < 1 {
		p.page = 1
	}
	return p.page, nil
}

func (p *Pager) Limit() int {
	return p.limit
}

func (p *Pager) Offset() int {
	p.offset = p.limit * (p.page - 1)
	return p.offset
}

func (p *Pager) Last() int {
	p.lastOnce.Do(func() { // TODO: 特にこうする意味もない気がする
		p.last = int(math.Ceil(float64(p.total) / float64(p.limit)))
	})
	return p.last
}

func (p *Pager) Prev() int {
	return p.page - 1
}

func (p *Pager) Next() int {
	return p.page + 1
}
