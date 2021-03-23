// gopaginator package helps in creating a single pagination object
package gopaginator

import (
	"log"
	"net/http"
	"strconv"
)

// default values for pagination
const (
	defaultLimit         int    = 10
	defaultPage          int    = 1
	defaultLimitParam    string = "limit"
	defaultPageParam     string = "page"
	defaultSearchParam   string = "q"
	defaultOrderByParam  string = "orderBy"
	defaultOrderingParam string = "order"
	OrderingASC          string = "ASC"
	OrderingDESC         string = "DESC"
)

// Paginator for pagination
type Paginator struct {
	request       *http.Request // for querying data from request
	Search        string        // to user in mysql Query, with % as boundary
	Limit         int           // for number of records
	Offset        int           // internally calculated as per page number and limit
	Ordering      string        // asc, desc
	OrderBy       string        // field_name
	Q             string        // to print the string sent by user
	page          int           // page number, used for calculating offset
	orderByParam  string        // orderBy param name in request
	orderingParam string        // ordering param name in request
	limitParam    string        // limit param name in request
	pageParam     string        // page param name in request
	qParam        string        // search query param name in request
}

// NewPaginator for creating a pagination object
func NewPaginator(r *http.Request) *Paginator {
	p := &Paginator{
		request:       r,
		limitParam:    defaultLimitParam,
		pageParam:     defaultPageParam,
		qParam:        defaultSearchParam,
		orderByParam:  defaultOrderByParam,
		orderingParam: defaultOrderingParam,
		Limit:         defaultLimit,
		page:          defaultPage,
		Offset:        0,
		Ordering:      OrderingASC,
		OrderBy:       "",
	}

	return p
}

// SetLimit for setting limit
func (p *Paginator) SetLimit(limit int, param string) *Paginator {
	if limit > 0 {
		p.Limit = limit
	}

	if param != "" {
		p.limitParam = param
	}

	return p
}

// SetPage for setting page
func (p *Paginator) SetPage(page int, param string) *Paginator {
	if page > 0 {
		p.page = page
	}
	if param != "" {
		p.pageParam = defaultPageParam
	}
	return p
}

// SetQ for setting search related params
func (p *Paginator) SetQ(q string, param string) *Paginator {
	p.Search = "%" + q + "%"
	p.Q = q
	if param != "" {
		p.qParam = defaultSearchParam
	}
	return p
}

// SetOrderBy for ordering and order by field
func (p *Paginator) SetOrderBy(orderBy string, param string) *Paginator {
	if orderBy != "" {
		p.OrderBy = orderBy
	}
	if param != "" {
		p.orderByParam = param
	}
	return p
}

// SetOrdering for ordering (ASC or DESC)
func (p *Paginator) SetOrdering(ordering string, param string) *Paginator {
	if ordering == OrderingASC || ordering == OrderingDESC {
		p.Ordering = ordering
	}
	if param != "" {
		p.orderingParam = param
	}
	return p
}

// SetOffset for calculating offset
func (p *Paginator) SetOffset() *Paginator {
	p.Offset = (p.page - 1) * p.Limit
	return p
}

// Parse for parsing query params from http.Request
func (p *Paginator) ParseRequest() *Paginator {
	values := p.request.URL.Query()

	// limit
	limit := values.Get(p.limitParam)
	limitInt, err := strconv.Atoi(limit)
	logError(err)

	// page
	page := values.Get(p.pageParam)
	pageInt, err := strconv.Atoi(page)
	logError(err)

	searchQuery := values.Get(p.qParam)
	order := values.Get(p.orderingParam)
	orderBy := values.Get(p.orderByParam)

	// create offset
	p.SetLimit(limitInt, p.limitParam).
		SetPage(pageInt, p.pageParam).
		SetQ(searchQuery, p.qParam).
		SetOrdering(order, p.orderingParam).
		SetOrderBy(orderBy, p.orderByParam).
		SetOffset()

	return p
}

func logError(err error) {
	if err != nil {
		log.Println("Error while parsing query params for pagination:", err.Error())
	}
}
