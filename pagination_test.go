package paginator_test

import (
	"net/http"
	paginator "saketsharma0805/pagination"
	"strconv"
	"testing"
)

func TestPagination(t *testing.T) {
	testCases := []struct {
		limit       int
		page        int
		orderBy     string
		ordering    string
		q           string
		expLimit    int
		expPage     int
		expOffset   int
		expSearch   string
		expQ        string
		expOrderBy  string
		expOrdering string
	}{
		{
			limit: 10, page: 1, orderBy: "name", ordering: "asc", q: "",
			expLimit: 10, expPage: 1, expOffset: 0, expSearch: "%%", expQ: "",
			expOrderBy: "name", expOrdering: "asc",
		},
		{
			limit: 10, page: 1, orderBy: "name", ordering: "desc", q: "",
			expLimit: 10, expPage: 1, expOffset: 0, expSearch: "%%", expQ: "",
			expOrderBy: "name", expOrdering: "desc",
		},
		{
			limit: 0, page: -1, orderBy: "na?.;me", ordering: "?abc", q: "abc;?123",
			expLimit: 10, expPage: 1, expOffset: 0, expSearch: "%abc123%", expQ: "abc;?123",
			expOrderBy: "name", expOrdering: "asc",
		},
	}

	for _, tc := range testCases {
		t.Run("Case", func(t *testing.T) {
			r, _ := http.NewRequest(http.MethodGet, "/", nil)

			// add query params
			q := r.URL.Query()
			q.Add("limit", strconv.Itoa(tc.limit))
			q.Add("page", strconv.Itoa(tc.page))
			q.Add("orderBy", tc.orderBy)
			q.Add("ordering", tc.ordering)
			q.Add("q", tc.q)
			r.URL.RawQuery = q.Encode()

			paginateObj := paginator.NewPagination(r, []string{"is_active"})
			paginateObj.Build()
			p := paginateObj.GetParams()

			if p.Limit != tc.expLimit {
				t.Errorf("Limit Expected: %d, Got: %d\n", tc.expLimit, p.Limit)
			}
			if p.Page != tc.expPage {
				t.Errorf("Page Expected: %d, Got: %d\n", tc.expPage, p.Page)
			}
			if p.Offset != tc.expOffset {
				t.Errorf("Offset Expected: %d, Got: %d\n", tc.expOffset, p.Offset)
			}
			if p.Q != tc.expQ {
				t.Errorf("Q Expected: %q, Got: %q\n", tc.expQ, p.Q)
			}
			if p.Search != tc.expSearch {
				t.Errorf("Search Expected: %q, Got: %q\n", tc.expSearch, p.Search)
			}
			if p.OrderBy != tc.expOrderBy {
				t.Errorf("OrderBy Expected: %q, Got: %q", tc.expOrderBy, p.OrderBy)
			}
			if p.Ordering != tc.expOrdering {
				t.Errorf("Ordering Expected: %q, Got: %q", tc.expOrdering, p.Ordering)
			}

		})
	}
}

func TestPagination2(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/", nil)

	// add query params
	q := r.URL.Query()
	q.Add("q", "shan")
	q.Add("is_active", "1")
	q.Add("email", "1")
	r.URL.RawQuery = q.Encode()

	paginateObj := paginator.NewPagination(r, []string{"is_active"})
	paginateObj.Build()
	p := paginateObj.GetParams()

	if p.Limit != 10 {
		t.Errorf("Limit Expected: %d, Got: %d\n", 10, p.Limit)
	}
	if p.Page != 1 {
		t.Errorf("Page Expected: %d, Got: %d\n", 1, p.Page)
	}
	if p.Offset != 0 {
		t.Errorf("Offset Expected: %d, Got: %d\n", 0, p.Offset)
	}
	if p.Q != "shan" {
		t.Errorf("Q Expected: %q, Got: %q\n", "shan", p.Q)
	}
	if p.Search != "%shan%" {
		t.Errorf("Search Expected: %q, Got: %q\n", "%shan%", p.Search)
	}
	if p.OrderBy != "" {
		t.Errorf("OrderBy Expected: %q, Got: %q", "", p.OrderBy)
	}
	if p.Ordering != "asc" {
		t.Errorf("Ordering Expected: %q, Got: %q", "asc", p.Ordering)
	}

	if p.Params["is_active"] != "1" {
		t.Errorf("Is active filter, Expected: %q, Got: %q", "1", p.Params["is_active"])
	}

	if p.Params["email"] != nil {
		t.Errorf("Unknown field, should be null")
	}

}
