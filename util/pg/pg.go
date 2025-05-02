package pg

import (
	"math"
	"sync"
)

// Deprecated: use pagination
type Pager struct {
	page     int
	lastOnce sync.Once
	last     int

	total int

	limit int
}

// Deprecated: use pagination
func New(total int) *Pager {
	return &Pager{
		total: total,
		limit: 10,
	}
}

// Deprecated: use pagination
func (p *Pager) Page(page int) int {
	if page < 1 {
		p.page = 1
		return p.page
	}
	p.page = page
	return p.page
}

// Deprecated: use pagination
func (p *Pager) Limit() int {
	return p.limit
}

// Deprecated: use pagination
func (p *Pager) Offset() int {
	return p.limit * (p.page - 1)
}

// Deprecated: use pagination
func (p *Pager) Last() int {
	p.lastOnce.Do(func() { // TODO: 特にこうする意味もない気がする
		p.last = int(math.Ceil(float64(p.total) / float64(p.limit)))
	})
	return p.last
}

// Deprecated: use pagination
func (p *Pager) Prev() int {
	return p.page - 1
}

// Deprecated: use pagination
func (p *Pager) Next() int {
	return p.page + 1
}
