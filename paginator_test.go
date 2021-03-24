package gopaginator_test

import (
	"net/http"
	"net/http/httptest"
	"saketsharma0805/gopaginator"
	"testing"
)

func TestLimitAndSearch(t *testing.T) {

	TestCases := []struct {
		Name           string
		Q              string
		Limit          string
		Page           string
		ExpectedSearch string
		ExpectedQ      string
		ExpectedLimit  int
		ExpectedOffset int
	}{
		{Name: "Case 1", Limit: "10", Page: "1", Q: "", ExpectedLimit: 10, ExpectedOffset: 0, ExpectedSearch: "%%", ExpectedQ: ""},
		{Name: "Case 2", Limit: "10", Page: "0", Q: "", ExpectedLimit: 10, ExpectedOffset: 0, ExpectedSearch: "%%", ExpectedQ: ""},
		{Name: "Case 3", Limit: "0", Page: "0", Q: "", ExpectedLimit: 10, ExpectedOffset: 0, ExpectedSearch: "%%", ExpectedQ: ""},
		{Name: "Case 4", Limit: "0", Page: "1", Q: "", ExpectedLimit: 10, ExpectedOffset: 0, ExpectedSearch: "%%", ExpectedQ: ""},
		{Name: "Case 5", Limit: "10", Page: "2", Q: "%", ExpectedLimit: 10, ExpectedOffset: 10, ExpectedSearch: "%%", ExpectedQ: "%"},
		{Name: "Case 6", Limit: "10", Page: "1", Q: "?", ExpectedLimit: 10, ExpectedOffset: 0, ExpectedSearch: "%%", ExpectedQ: "?"},
		{Name: "Case 7", Limit: "10", Page: "1", Q: "%", ExpectedLimit: 10, ExpectedOffset: 0, ExpectedSearch: "%%", ExpectedQ: "%"},
		{Name: "Case 8", Limit: "10", Page: "1", Q: "paginator", ExpectedLimit: 10, ExpectedOffset: 0, ExpectedSearch: "%paginator%", ExpectedQ: "paginator"},
		{Name: "Case 9", Limit: "10", Page: "1", Q: "paginator;", ExpectedLimit: 10, ExpectedOffset: 0, ExpectedSearch: "%paginator%", ExpectedQ: "paginator;"},
		{Name: "Case 10", Limit: "10", Page: "1", Q: "testing.server@gmail.com", ExpectedLimit: 10, ExpectedOffset: 0, ExpectedSearch: "%testing.server@gmail.com%", ExpectedQ: "testing.server@gmail.com"},
	}

	for _, tc := range TestCases {
		t.Run(tc.Name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			q := r.URL.Query()
			q.Add(gopaginator.DefaultLimitParam, tc.Limit)
			q.Add(gopaginator.DefaultPageParam, tc.Page)
			q.Add(gopaginator.DefaultSearchParam, tc.Q)
			r.URL.RawQuery = q.Encode()

			p := gopaginator.NewPaginator(r, []string{})
			p.ParseRequest()

			if tc.ExpectedLimit != p.Limit {
				t.Errorf("Expected limit %d, Got %d", tc.ExpectedLimit, p.Limit)
			}

			if tc.ExpectedOffset != p.Offset {
				t.Errorf("Expected offset %d, Got %d", tc.ExpectedOffset, p.Offset)
			}

			if tc.ExpectedQ != p.Q {
				t.Errorf("Expected q param %q, Got %q", tc.ExpectedQ, p.Q)
			}

			if tc.ExpectedSearch != p.Search {
				t.Errorf("Expected search query %q, Got %q", tc.ExpectedSearch, p.Search)
			}
		})
	}
}

func TestOrder(t *testing.T) {

	TestCases := []struct {
		Name            string
		Order           string
		OrderBy         string
		ExpectedOrder   string
		ExpectedOrderBy string
	}{
		{Name: "Case 1", Order: "", OrderBy: "first_name", ExpectedOrder: gopaginator.OrderingASC, ExpectedOrderBy: "first_name"},
		{Name: "Case 2", Order: "desc", OrderBy: "first_name", ExpectedOrder: gopaginator.OrderingDESC, ExpectedOrderBy: "first_name"},
		{Name: "Case 3", Order: "asc", OrderBy: "first_?nam-e", ExpectedOrder: gopaginator.OrderingASC, ExpectedOrderBy: "first_name"},
	}

	for _, tc := range TestCases {
		t.Run(tc.Name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			q := r.URL.Query()
			q.Add(gopaginator.DefaultOrderingParam, tc.Order)
			q.Add(gopaginator.DefaultOrderByParam, tc.OrderBy)
			r.URL.RawQuery = q.Encode()

			p := gopaginator.NewPaginator(r, []string{})
			p.ParseRequest()

			if tc.ExpectedOrder != p.Ordering {
				t.Errorf("Expected order %q, Got %q", tc.ExpectedOrder, p.Ordering)
			}

			if tc.ExpectedOrderBy != p.OrderBy {
				t.Errorf("Expected orderBy %q, Got %q", tc.ExpectedOrderBy, p.OrderBy)
			}
		})
	}
}

func TestParamNames(t *testing.T) {

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	q := r.URL.Query()
	q.Add("l", "50")
	q.Add("p", "2")
	q.Add("query", "s")
	q.Add("o", "desc")
	q.Add("is_active", "1")
	q.Add("oby", "")
	r.URL.RawQuery = q.Encode()

	p := gopaginator.NewPaginator(r, []string{"is_active"})
	p.SetLimit("l", 10).
		SetPage("p", 1).
		SetOrderBy("oby", "email").
		SetOrdering("o", "asc").
		SetQ("query", "").
		ParseRequest()

	if p.Limit != 50 {
		t.Errorf("Error while verifying limit, Expected %d, Got %d", 50, p.Limit)
	}

	if p.Offset != 50 {
		t.Errorf("Error while verifying offset, Expected %d, Got %d", 50, p.Offset)
	}

	if p.Search != "%s%" {
		t.Errorf("Error while verifying search query, Expected %q, Got %q", "%s%", p.Search)
	}

	if p.Ordering != "DESC" {
		t.Errorf("Error while verifying order, Expected %q, Got %q", "DESC", p.Ordering)
	}

	if p.OrderBy != "email" {
		t.Errorf("Error while verifying orderBy field, Expected %q, Got %q", "email", p.OrderBy)
	}

	if val, ok := p.Filter["is_active"]; !ok || val != "1" {
		t.Errorf("Expected is_active value to be %q, Got %q", "1", val)
	}
}
