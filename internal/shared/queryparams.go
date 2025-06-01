package shared

import (
	"github.com/google/go-querystring/query"
)

type QueryParams struct {
	Page  int    `url:"page"`
	S     string `url:"s,omitempty"`
	Color string `url:"color,omitempty"`
}

func (q *QueryParams) SetPage(p int) *QueryParams {
	q.Page = p
	return q
}

func (q *QueryParams) GetQuery() string {
	v, _ := query.Values(q)
	return "?" + v.Encode()
}
