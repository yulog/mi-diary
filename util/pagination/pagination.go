package pagination

import (
	"fmt"
	"log/slog"
	"math"
)

type Pages struct {
	Current int
	Prev    Page
	Next    Page
	Last    Page
}

type Page struct {
	Index int
	Has   bool
}

type NextPageChecker interface {
	HasNextPage(p *Pagination) bool
}

type PreviousPageChecker interface {
	HasPreviousPage(p *Pagination) bool
}

type LastPageChecker interface {
	HasLastPage(p *Pagination) bool
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
	LastChecker     LastPageChecker     // 最終ページ判定ロジック
}

type DefaultNextPageChecker struct{}

func (c DefaultNextPageChecker) HasNextPage(p *Pagination) bool {
	return p.CurrentPage < p.TotalPages()
}

type DefaultPreviousPageChecker struct{}

func (c DefaultPreviousPageChecker) HasPreviousPage(p *Pagination) bool {
	return p.CurrentPage > 1
}

type DefaultLastPageChecker struct{}

func (c DefaultLastPageChecker) HasLastPage(p *Pagination) bool {
	return false
}

func New(currentPage, perPage, totalItems int, nchecker NextPageChecker, pchecker PreviousPageChecker, lchecker LastPageChecker) (*Pagination, error) {
	if nchecker == nil {
		nchecker = DefaultNextPageChecker{} // デフォルトロジックを使用
	}
	if pchecker == nil {
		pchecker = DefaultPreviousPageChecker{}
	}
	if lchecker == nil {
		lchecker = DefaultLastPageChecker{}
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

func (p *Pagination) HasLastPage() bool {
	return p.LastChecker.HasLastPage(p)
}

// TODO: 0ではなくCurrentPageを返しても良いかも？
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

func (p *Pagination) Pages() Pages {
	next, err := p.NextPage()
	if err != nil {
		slog.Info(err.Error())
	}
	prev, err := p.PreviousPage()
	if err != nil {
		slog.Info(err.Error())
	}
	return Pages{
		Current: p.CurrentPage,
		Prev:    Page{Index: prev, Has: p.HasPreviousPage()},
		Next:    Page{Index: next, Has: p.HasNextPage()},
		Last:    Page{Index: p.TotalPages(), Has: p.HasLastPage()},
	}
}
