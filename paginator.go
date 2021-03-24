// gopaginator package helps in creating a single pagination object
package gopaginator

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// default values for pagination
const (
	DefaultLimitParam    string = "limit"
	DefaultPageParam     string = "page"
	DefaultSearchParam   string = "q"
	DefaultOrderByParam  string = "orderBy"
	DefaultOrderingParam string = "order"
	OrderingASC          string = "ASC"
	OrderingDESC         string = "DESC"
	sanitizeSearch       string = "[^a-zA-Z0-9\\.\\_\\-@]"
	sanitizeField        string = "[^a-zA-Z\\_]"
	defaultLimit         int    = 10
	defaultPage          int    = 1
)

// sanitize search query and database field names
var regSearch = regexp.MustCompile(sanitizeSearch)
var regField = regexp.MustCompile(sanitizeField)

// Paginator for pagination
type Paginator struct {
	request       *http.Request          // for querying data from request
	Search        string                 // to user in mysql Query, with % as boundary
	Limit         int                    // for number of records
	Offset        int                    // internally calculated as per page number and limit
	Ordering      string                 // asc, desc
	OrderBy       string                 // field_name
	Q             string                 // to print the string sent by user, [do not use in search db]
	Filter        map[string]interface{} // for saving query param's key and value
	page          int                    // page number, used for calculating offset
	orderByParam  string                 // orderBy param name in request
	orderingParam string                 // ordering param name in request
	limitParam    string                 // limit param name in request
	pageParam     string                 // page param name in request
	qParam        string                 // search query param name in request
	filters       []string
}

// NewPaginator for creating a pagination object
func NewPaginator(r *http.Request, extraParams []string) *Paginator {
	p := &Paginator{
		request:       r,
		limitParam:    DefaultLimitParam,
		pageParam:     DefaultPageParam,
		qParam:        DefaultSearchParam,
		orderByParam:  DefaultOrderByParam,
		orderingParam: DefaultOrderingParam,
		Limit:         defaultLimit,
		page:          defaultPage,
		Offset:        0,
		Ordering:      OrderingASC,
		OrderBy:       "",
		Filter:        make(map[string]interface{}),
		filters:       extraParams,
	}

	return p
}

// SetLimit for setting limit
func (p *Paginator) SetLimit(param string, limit int) *Paginator {
	if limit > 0 {
		p.Limit = limit
	}

	if param != "" {
		p.limitParam = param
	}

	return p
}

// SetPage for setting page
func (p *Paginator) SetPage(param string, page int) *Paginator {
	if page > 0 {
		p.page = page
	}
	if param != "" {
		p.pageParam = param
	}
	return p
}

// SetQ for setting search related params
func (p *Paginator) SetQ(param string, q string) *Paginator {
	res := regSearch.ReplaceAllString(q, "")
	p.Search = "%" + res + "%"
	p.Q = q
	if param != "" {
		p.qParam = param
	}
	return p
}

// SetOrderBy for ordering and order by field
func (p *Paginator) SetOrderBy(param string, orderBy string) *Paginator {
	orderBy = regField.ReplaceAllString(orderBy, "${1}")
	if orderBy != "" {
		p.OrderBy = orderBy
	}
	if param != "" {
		p.orderByParam = param
	}
	return p
}

// SetOrdering for ordering (ASC or DESC)
func (p *Paginator) SetOrdering(param string, ordering string) *Paginator {
	ordering = strings.ToUpper(ordering)
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
	if err == nil {
		p.SetLimit(p.limitParam, limitInt)
	}

	// page
	page := values.Get(p.pageParam)
	pageInt, err := strconv.Atoi(page)
	if err == nil {
		p.SetPage(p.pageParam, pageInt)
	}

	searchQuery := values.Get(p.qParam)
	order := values.Get(p.orderingParam)
	orderBy := values.Get(p.orderByParam)

	// create offset
	p.SetQ(p.qParam, searchQuery).
		SetOrdering(p.orderingParam, order).
		SetOrderBy(p.orderByParam, orderBy).
		SetOffset()

	for _, val := range p.filters {
		qValue := values.Get(val)
		p.Filter[val] = qValue
	}

	return p
}

// func logError(err error) {
// 	if err != nil {
// 		log.Println("Error while parsing query params for pagination:", err.Error())
// 	}
// }
