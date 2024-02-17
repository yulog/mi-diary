package pg

import (
	"math"
	"sync"
)

type Pager struct {
	page int
	// Prev int
	// Next int
	lastOnce sync.Once
	last     int

	total int

	limit  int
	offset int
}

func New(total int) *Pager {
	return &Pager{
		// c: c,

		total: total,
		limit: 10,
	}
}

func (p *Pager) Page(page int) int {
	p.page = page
	if p.page < 1 {
		p.page = 1
	}
	return p.page
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
