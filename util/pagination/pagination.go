package pagination

import (
	"fmt"
	"math"
)

type NextPageChecker interface {
	HasNextPage(p *Pagination) bool
}

type PreviousPageChecker interface {
	HasPreviousPage(cp *Pagination) bool
}

type Paging interface {
	Offset() int
	Limit() int
}

type Pagination struct {
	CurrentPage     int
	PerPage         int // 1ページあたりの項目数
	TotalItems      int
	NextChecker     NextPageChecker     // 次ページ判定ロジック
	PreviousChecker PreviousPageChecker // 前ページ判定ロジック
}

type DefaultNextPageChecker struct{}

func (c DefaultNextPageChecker) HasNextPage(p *Pagination) bool {
	return p.CurrentPage < p.TotalPages()
}

type DefaultPreviousPageChecker struct{}

func (c DefaultPreviousPageChecker) HasPreviousPage(p *Pagination) bool {
	return p.CurrentPage > 1
}

func New(currentPage, perPage, totalItems int, nchecker NextPageChecker, pchecker PreviousPageChecker) (*Pagination, error) {
	if nchecker == nil {
		nchecker = DefaultNextPageChecker{} // デフォルトロジックを使用
	}
	if pchecker == nil {
		pchecker = DefaultPreviousPageChecker{}
	}
	if currentPage <= 0 {
		currentPage = 1
		// return nil, fmt.Errorf("current page must be positive")
	}
	if perPage <= 0 {
		return nil, fmt.Errorf("per page must be positive")
	}
	if totalItems < 0 {
		return nil, fmt.Errorf("total items must be non-negative")
	}
	return &Pagination{
		CurrentPage:     currentPage,
		PerPage:         perPage,
		TotalItems:      totalItems,
		NextChecker:     nchecker,
		PreviousChecker: pchecker,
	}, nil
}

func (p *Pagination) TotalPages() int {
	if p.PerPage == 0 {
		return 0
	}
	return int(math.Ceil(float64(p.TotalItems) / float64(p.PerPage)))
}

func (p *Pagination) Offset() int {
	return (p.CurrentPage - 1) * p.PerPage
}

func (p *Pagination) Limit() int {
	return p.PerPage
}

func (p *Pagination) HasNextPage() bool {
	return p.NextChecker.HasNextPage(p)
}

func (p *Pagination) HasPreviousPage() bool {
	return p.PreviousChecker.HasPreviousPage(p)
}

func (p *Pagination) NextPage() (int, error) {
	if !p.HasNextPage() {
		return 0, fmt.Errorf("no next page")
	}
	return p.CurrentPage + 1, nil
}

func (p *Pagination) PreviousPage() (int, error) {
	if !p.HasPreviousPage() {
		return 0, fmt.Errorf("no previous page")
	}
	return p.CurrentPage - 1, nil
}
