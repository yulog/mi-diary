package pagination

import (
	"testing"
)

type AlwaysHasNextPageChecker struct{}

func (c AlwaysHasNextPageChecker) HasNextPage(p *Pagination) bool {
	return p.CurrentPage < p.TotalPages()+1
}

type ComplexChecker struct{}

func (c ComplexChecker) HasNextPage(p *Pagination) bool {
	if p.TotalPages() > 10 && p.CurrentPage <= 5 {
		return true
	}
	return p.CurrentPage < p.TotalPages()
}

func TestHasNextPage(t *testing.T) {
	tests := []struct {
		name        string
		currentPage int
		perPage     int
		totalItems  int
		checker     NextPageChecker
		want        bool
	}{
		{
			name:        "デフォルト: 次ページあり",
			currentPage: 1,
			perPage:     10,
			totalItems:  50,
			checker:     nil, // デフォルトチェッカー
			want:        true,
		},
		{
			name:        "デフォルト: 次ページなし",
			currentPage: 5,
			perPage:     10,
			totalItems:  50,
			checker:     nil,
			want:        false,
		},
		{
			name:        "デフォルト: TotalItemsが0でCurrentPageが1の場合",
			currentPage: 1,
			perPage:     10,
			totalItems:  0,
			want:        false,
		},
		{
			name:        "デフォルト: CurrentPageがTotalPagesと同じ場合",
			currentPage: 3,
			perPage:     10,
			totalItems:  30,
			want:        false,
		},
		{
			name:        "デフォルト: CurrentPageがTotalPagesより大きい場合",
			currentPage: 4,
			perPage:     10,
			totalItems:  30,
			want:        false,
		},
		{
			name:        "カスタム: 常に次ページあり",
			currentPage: 5,
			perPage:     10,
			totalItems:  50,
			checker:     AlwaysHasNextPageChecker{},
			want:        true,
		},
		{
			name:        "カスタム: 複雑な条件1",
			currentPage: 3,
			perPage:     10,
			totalItems:  150,
			checker:     ComplexChecker{},
			want:        true,
		},
		{
			name:        "カスタム: 複雑な条件2",
			currentPage: 7,
			perPage:     10,
			totalItems:  150,
			checker:     ComplexChecker{},
			want:        true,
		},
		{
			name:        "カスタム: 複雑な条件3",
			currentPage: 12,
			perPage:     10,
			totalItems:  150,
			checker:     ComplexChecker{},
			want:        true,
		},
		{
			name:        "カスタム: 複雑な条件4",
			currentPage: 15,
			perPage:     10,
			totalItems:  150,
			checker:     ComplexChecker{},
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, _ := New(tt.currentPage, tt.perPage, tt.totalItems, tt.checker, nil)
			if got := p.HasNextPage(); got != tt.want {
				t.Errorf("HasNextPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNextPage(t *testing.T) {
	tests := []struct {
		name        string
		currentPage int
		perPage     int
		totalItems  int
		want        int
		wantErr     bool
	}{
		{
			name:        "次のページあり",
			currentPage: 1,
			perPage:     10,
			totalItems:  25,
			want:        2,
			wantErr:     false,
		},
		{
			name:        "次のページなし",
			currentPage: 3,
			perPage:     10,
			totalItems:  25,
			want:        0,
			wantErr:     true,
		},
		{
			name:        "TotalItemsが0の場合",
			currentPage: 1,
			perPage:     10,
			totalItems:  0,
			want:        0,
			wantErr:     true, // TotalPagesが0になるため次ページは存在しない
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, _ := New(tt.currentPage, tt.perPage, tt.totalItems, nil, nil)
			got, err := p.NextPage()
			if (err != nil) != tt.wantErr {
				t.Errorf("NextPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NextPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPreviousPage(t *testing.T) {
	tests := []struct {
		name        string
		currentPage int
		perPage     int
		totalItems  int
		want        int
		wantErr     bool
	}{
		{
			name:        "前のページあり",
			currentPage: 2,
			perPage:     10,
			totalItems:  25,
			want:        1,
			wantErr:     false,
		},
		{
			name:        "前のページなし",
			currentPage: 1,
			perPage:     10,
			totalItems:  25,
			want:        0,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, _ := New(tt.currentPage, tt.perPage, tt.totalItems, nil, nil)
			got, err := p.PreviousPage()
			if (err != nil) != tt.wantErr {
				t.Errorf("PreviousPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PreviousPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name            string
		currentPage     int
		perPage         int
		totalItems      int
		wantCurrentPage int
		wantErr         bool
	}{
		{
			name:            "正常",
			currentPage:     1,
			perPage:         10,
			totalItems:      100,
			wantCurrentPage: 1,
			wantErr:         false,
		},
		{
			name:            "currentPageが0",
			currentPage:     0,
			perPage:         10,
			totalItems:      100,
			wantCurrentPage: 1,
			wantErr:         false, // エラーではなく補正される
		},
		{
			name:            "currentPageが負の値",
			currentPage:     -1,
			perPage:         10,
			totalItems:      100,
			wantCurrentPage: 1,
			wantErr:         false, // エラーではなく補正される
		},
		{
			name:        "perPageが0",
			currentPage: 1,
			perPage:     0,
			totalItems:  100,
			wantErr:     true,
		},
		{
			name:        "totalItemsが負の値",
			currentPage: 1,
			perPage:     10,
			totalItems:  -1,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := New(tt.currentPage, tt.perPage, tt.totalItems, nil, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && p.CurrentPage != tt.wantCurrentPage {
				t.Errorf("New() currentPage = %d, wantCurrentPage %d", p.CurrentPage, tt.wantCurrentPage)
			}
		})
	}
}
