package paginator

import (
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const (
	OrderingASC  string = "asc"
	OrderingDESC string = "desc"
)

var (
	Limit          int    = 10
	Page           int    = 1
	Offset         int    = 0
	Q              string = ""
	OrderBy        string = ""
	sanitizeSearch string = "[^a-zA-Z0-9\\.\\_\\-@]"
	sanitizeField  string = "[^a-zA-Z\\_]"
)

var (
	regSearch = regexp.MustCompile(sanitizeSearch)
	regField  = regexp.MustCompile(sanitizeField)
)

type Paginator interface {
	SetOrderBy(OrderBy string)
	SetOrdering(Ordering string)
	SetQuery(Q string)
	SetPage(Page int)
	SetLimit(Limit int)
	Build()
	GetParams() *Pagination
}

// Pagination for pagination object
type Pagination struct {
	err      error
	request  *http.Request
	Limit    int
	Page     int
	Offset   int
	Q        string
	Search   string
	OrderBy  string
	Ordering string
	Filters  []string
	Params   map[string]interface{}
}

func (p *Pagination) GetParams() *Pagination {
	val := p.request.URL.Query()
	for _, v := range p.Filters {
		p.Params[v] = val.Get(v)
	}
	return p
}

func (p *Pagination) SetLimit(Limit int) {
	if Limit < 10 {
		Limit = 10
	}
	p.Limit = Limit
}

func (p *Pagination) SetPage(Page int) {
	if Page < 1 {
		Page = 1
	}
	p.Page = Page
	p.Offset = (p.Page - 1) * p.Limit
}

func (p *Pagination) SetQuery(Q string) {
	res := regSearch.ReplaceAllString(Q, "")
	p.Search = "%" + res + "%"
	p.Q = Q
}

func (p *Pagination) SetOrderBy(OrderBy string) {
	OrderBy = regField.ReplaceAllString(OrderBy, "${1}")
	p.OrderBy = OrderBy
}

func (p *Pagination) SetOrdering(ordering string) {
	ordering = strings.ToLower(ordering)
	if ordering == OrderingASC || ordering == OrderingDESC {
		p.Ordering = ordering
	}
}

func (p *Pagination) Build() {
	if p.request == nil {
		p.err = nil
		return
	}
	v := p.request.URL.Query()

	Limit := v.Get("limit")
	LimitInt, _ := strconv.Atoi(Limit)
	p.SetLimit(LimitInt)

	Page := v.Get("page")
	PageInt, _ := strconv.Atoi(Page)
	p.SetPage(PageInt)

	Query := v.Get("q")
	p.SetQuery(Query)

	OrderBy := v.Get("orderBy")
	p.SetOrderBy(OrderBy)

	Ordering := v.Get("ordering")
	p.SetOrdering(Ordering)

}

// NewPagination for creating pagination instance with default values
func NewPagination(r *http.Request, filters []string) Paginator {
	return &Pagination{
		err:      errors.New("call build method for parsing reQuest"),
		request:  r,
		Limit:    Limit,
		Page:     Page,
		Offset:   Offset,
		Q:        Q,
		OrderBy:  OrderBy,
		Ordering: OrderingASC,
		Search:   Q,
		Filters:  filters,
		Params:   make(map[string]interface{}),
	}
}
